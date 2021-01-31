package main

import (
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	dtypes "github.com/Prev/HotFunctions/worker_front/types"
)

type FunctionRunner struct {
	imageBuilder               *ImageBuilder
	singletonContainerManager  *SingletonContainerManager
	cachingOptions             CachingOptions
	lru                        map[string]int64
	images                     map[string]Image
	containers                 []Container
	mutex                      *sync.Mutex
}

func newFunctionRunner(cachingOptions CachingOptions) *FunctionRunner {
	r := new(FunctionRunner)
	r.cachingOptions = cachingOptions
	r.imageBuilder = newImageBuilder(&r.cachingOptions)
	r.singletonContainerManager = newRestContainerManager()

	r.lru = make(map[string]int64)
	r.images = make(map[string]Image)
	r.containers = make([]Container, 0)
	r.mutex = new(sync.Mutex)
	return r
}

// Prepare all images of sample functions
// [Notice] It may takes long time, so call it carefully
func (r *FunctionRunner) PrepareImages() error {
	resp, err := http.Get(UserFunctionUrlPrefix + "lists.txt")
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	lists := strings.Split(string(data), "\n")
	logger.Println("Build image", strings.Join(lists, ","))

	for _, functionName := range lists {
		if functionName == "" {
			continue
		}
		image, err := r.imageBuilder.BuildSafe(functionName)
		if err != nil {
			logger.Println(err.Error())
		}
		r.images[functionName] = image
	}
	return nil
}

func (r *FunctionRunner) Reset(resetImages bool) {
	CleanContainers()
	r.lru = make(map[string]int64)
	r.containers = make([]Container, 0)
	r.singletonContainerManager = newRestContainerManager()

	if resetImages {
		r.images = make(map[string]Image)
	}
}

func (r *FunctionRunner) RunFunction(functionName string) (*dtypes.ContainerResponse, dtypes.FunctionExecutionMetaData) {
	var err error
	tryCnt := 0
	meta := dtypes.FunctionExecutionMetaData{
		ImageBuilt: false,
		UsingPooledContainer: false,
		UsingExistingRestContainer: false,
		ContainerName: "",
		ImageName: "",
	}

	// Step1: Check for the image existence.
	// If there is no cached image, build a new docker image
	r.mutex.Lock()
	image, imageExists := r.images[functionName]
	// Set lru as maximum value while running
	r.lru[functionName] = math.MaxInt64
	r.mutex.Unlock()

BuildImage:
	if imageExists == false {
		// If there is no image, create a new one
		if image, err = r.imageBuilder.BuildSafe(functionName); err != nil {
			logger.Println(err.Error())
		}
		meta.ImageBuilt = true
		r.mutex.Lock()
		r.images[functionName] = image
		r.mutex.Unlock()
	}

	// Step2: Check for the container existence.
	// There are two cases of using existence containers: 1) pre-warmed or 2) restful
	// Rest containers are capable to execute multiple functions concurrently in single container,
	// so it requires singleton pattern, where `r.singletonContainerManager` is in charge for it.

SelectContainer:
	selected := Container{}

	if r.cachingOptions.UsingRestMode {
		// RestMode
		if cont, exists := r.singletonContainerManager.Get(image); exists {
			selected = cont
			meta.UsingExistingRestContainer = true

		} else {
			selected, err = r.singletonContainerManager.SafeCreate(image)
			if err != nil {
				return nil, meta
			}
		}

	} else {
		// Pre-warmed container pool
		r.mutex.Lock()
		for i, cont := range r.containers {
			if cont.FunctionName == functionName {
				selected = cont
				// Since container is not reusable, remove it from the list
				r.containers = append(r.containers[:i], r.containers[i+1:]...)
				meta.UsingPooledContainer = true
				break
			}
		}
		r.mutex.Unlock()
	}

	if selected.Name == "" {
		// If there is no container can be used directly, create a new container
		selected, err = CreateContainer(image)
		if err != nil {
			logger.Println(err.Error())
			return nil, meta
		}
	}

	out, err := selected.Run()
	if err != nil {
		tryCnt++
		if tryCnt > 50 {
			return nil, meta
		}

		// Retry the process from the image building or selecting the container
		logger.Println("Error:", err.Error(), "Retry...", tryCnt)

		if selected.IsRestMode {
			time.Sleep(100 * time.Millisecond)
			goto SelectContainer

		} else {
			time.Sleep(100 * time.Millisecond)
			imageExists = false
			goto BuildImage
		}
	}

	meta.ContainerName = selected.Name
	meta.ImageName = image.Name

	r.mutex.Lock()
	r.lru[functionName] = time.Now().Unix()
	r.mutex.Unlock()

	go r.manageCaches()
	return out, meta
}

func (r *FunctionRunner) manageCaches() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check for the user code size limit
	// If the limit is set -1, the is no limitation
	if r.cachingOptions.UserCodeSizeLimit != -1 {
		totalCodeSize := int64(0)

		// Get total downloaded code size
		for _, functionName := range sortMapByValue(&r.lru, true) {
			if img, exists := r.images[functionName]; exists {
				totalCodeSize += img.Size
			}
		}

		// Remove from older images
		for _, functionName := range sortMapByValue(&r.lru, true) {
			if img, exists := r.images[functionName]; exists {
				if totalCodeSize > r.cachingOptions.UserCodeSizeLimit {
					if err := r.imageBuilder.RemoveImage(functionName); err != nil {
						logger.Println(err.Error())
					}

					logger.Printf("Image %s is removed (%.1fMB, limit: %.1fMB)\n",
						img.Name,
						float64(img.Size) / 1000000,
						float64(r.cachingOptions.UserCodeSizeLimit) / 1000000,
					)
					delete(r.images, functionName)
					totalCodeSize -= img.Size
				}
			}
		}
	}

	// Check for the container pooling
	// If the limit is set -1, do not use the feature
	if r.cachingOptions.ContainerPoolLimit != -1 {
		for i, functionName := range sortMapByValue(&r.lru, false) {
			image := r.images[functionName]
			if i < r.cachingOptions.ContainerPoolLimit {
				createPreWarmedContainers(&r.containers, image, r.cachingOptions.ContainerPoolNum)

			} else {
				clearPreWarmedContainers(&r.containers, image)
			}
		}
	}

	// Remove rest container based on lifetime
	if r.cachingOptions.UsingRestMode {
		now := time.Now().Unix()
		lifetime := int64(r.cachingOptions.RestContainerLifeTime)

		for _, container := range r.singletonContainerManager.containers {
			if now - r.lru[container.FunctionName] > lifetime {
				logger.Printf("Container %s is removed\n", container.Name)
				container.Remove()
				r.singletonContainerManager.Delete(container.FunctionName)
			}
		}
	}
}

