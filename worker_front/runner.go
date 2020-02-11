package main

import (
	"sync"
	"time"
)

type FunctionRunner struct {
	imageBuilder         *ImageBuilder
	cachingOptions       CachingOptions
	lru                  map[string]int64
	images               map[string]Image
	containers           []Container
	mutex                *sync.Mutex
}

type runningMetaData struct {
	ImageBuilt              bool
	UsingPooledContainer    bool
	UsingRestContainer      bool
	ContainerName           string
	ImageName               string
}

func newFunctionRunner(cachingOptions CachingOptions) *FunctionRunner {
	r := new(FunctionRunner)
	r.cachingOptions = cachingOptions
	r.imageBuilder = newImageBuilder(&r.cachingOptions)

	r.lru = make(map[string]int64)
	r.images = make(map[string]Image)
	r.mutex = new(sync.Mutex)
	return r
}

func (r *FunctionRunner) runFunction(functionName string) (*FunctionResponse, runningMetaData) {
	var err error
	meta := runningMetaData{false, false, false, "", ""}

	// Step1: Check for the image existence.
	// If there is no cached image, build a new docker image
	r.mutex.Lock()
	image, imageExists := r.images[functionName]
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
	// Pre-warmed containers are not reusable. They are created in background after function executions.

	// Rest containers are capable to execute multiple functions concurrently in single container,
	// so it means they are reusable, and the runner have to handle pool of containers by this "re-usability" feature
	r.mutex.Lock()
	selected := Container{}
	for i, cont := range r.containers {
		if cont.FunctionName == functionName {
			selected = cont
			if !selected.Reusable {
				// If the container is not reusable, remove it from the list
				r.containers = append(r.containers[:i], r.containers[i+1:]...)
			}
			// Record meta data
			if selected.IsRestMode {
				meta.UsingRestContainer = true
			} else {
				meta.UsingPooledContainer = true
			}
			break
		}
	}
	r.mutex.Unlock()

	if selected.Name == "" {
		// If there is no container can be used directly, create a new container
		selected, err = CreateContainer(image)
		if err != nil {
			logger.Println(err.Error())
			return nil, meta
		}

		if selected.Reusable {
			// If container is reusable, register to the container pool.
			r.mutex.Lock()
			r.containers = append(r.containers, selected)
			r.mutex.Unlock()
		}
	}

	out, err := selected.Run()
	if err != nil {
		// Retry the process from the image building.
		logger.Println(err.Error(), "Retry...")
		imageExists = false
		goto BuildImage
	}

	meta.ContainerName = selected.Name
	meta.ImageName = image.Name

	r.lru[functionName] = time.Now().Unix()
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
	}
}

