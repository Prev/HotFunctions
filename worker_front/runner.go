package main

import (
	"sync"
	"time"

	"github.com/Prev/HotFunctions/worker_front/types"
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
	r.mutex = new(sync.Mutex)
	return r
}

func (r *FunctionRunner) runFunction(functionName string) (*types.ContainerResponse, types.FunctionExecutionMetaData) {
	var err error
	tryCnt := 0
	meta := types.FunctionExecutionMetaData{false, false, false, "", ""}

	// Step1: Check for the image existence.
	// If there is no cached image, build a new docker image
	r.mutex.Lock()
	image, imageExists := r.images[functionName]
	r.lru[functionName] = time.Now().Unix()
	r.mutex.Unlock()

BuildImage:
	if imageExists == false {
		// If there is no image, create a new one
		if image, err = r.imageBuilder.BuildSafe(functionName); err != nil {
			logger.Println(err.Error())
		}
		meta.ImageBuilt = true
		r.images[functionName] = image
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
		if tryCnt > 30 {
			return nil, meta
		}

		// Retry the process from the image building or selecting the container
		logger.Println("Error:", err.Error(), "Retry...", tryCnt)

		if selected.IsRestMode {
			//if strings.Contains(err.Error(), "tcp") {
			//	selected.Remove()
			//	goto SelectContainer
			//} else {
			//	time.Sleep(100 * time.Millisecond)
			//	goto RunContainer
			//}
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

	for i, functionName := range sortMapByValue(&r.lru) {
		// Check for the image cache limit
		// If the limit is set -1, the is no limitation
		if r.cachingOptions.ImageLimit != -1 && i >= r.cachingOptions.ImageLimit {
			if _, exists := r.images[functionName]; exists {
				if err := r.imageBuilder.RemoveImage(functionName); err != nil {
					logger.Println(err.Error())
				}
				delete(r.images, functionName)
			}
		}

		// Check for the container pooling
		// If the limit is set -1, do not use the feature
		if r.cachingOptions.ContainerPoolLimit != -1 {
			image := r.images[functionName]
			if i < r.cachingOptions.ContainerPoolLimit {
				createPreWarmedContainers(&r.containers, image, r.cachingOptions.ContainerPoolNum)

			} else {
				clearPreWarmedContainers(&r.containers, image)
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
}

