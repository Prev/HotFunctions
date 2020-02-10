package main

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

type FunctionRunner struct {
	cachingOptions       CachingOptions
	lru                  map[string]int64
	imageExists          map[string]bool
	prewarmedContainers  []string
	imageBuilder         *ImageBuilder
	mutex                *sync.Mutex
}

type runningMetaData struct {
	ImageBuilt       bool
	ContainerCreated bool
}

func newFunctionRunner(cachingOptions CachingOptions) *FunctionRunner {
	r := new(FunctionRunner)
	r.cachingOptions = cachingOptions
	r.lru = make(map[string]int64)
	r.imageExists = make(map[string]bool)
	r.imageBuilder = newImageBuilder()
	r.mutex = new(sync.Mutex)
	return r
}

func (r *FunctionRunner) runFunction(functionName string) (*FunctionResponse, runningMetaData) {
	var err error
	meta := runningMetaData{false, false}

	// Step1: Check for the image existence.
	// If there is no cached image, build a new docker image
	r.mutex.Lock()
	imageExists := r.imageExists[functionName]
	r.mutex.Unlock()

BuildImage:
	if imageExists == false {
		// If there is no image, create a new one
		meta.ImageBuilt = true
		if err = r.imageBuilder.BuildSafe(functionName); err != nil {
			logger.Println(err.Error())
		}
	}

	// Step2: Check for the container existence.
	targetContainerName := ""
	r.mutex.Lock()
	for i, containerName := range r.prewarmedContainers {
		if containerBelongsToFunction(containerName, functionName) {
			// Pick the proper idle container and remove it from the list
			targetContainerName = containerName
			r.prewarmedContainers = append(r.prewarmedContainers[:i], r.prewarmedContainers[i+1:]...)
			break
		}
	}
	r.mutex.Unlock()

	if targetContainerName == "" {
		// If there is no available container, create a new one
		meta.ContainerCreated = true
		if targetContainerName, err = CreateContainer(functionName); err != nil {
			logger.Println(err.Error())
			return nil, meta
		}
	}

	out, err := RunContainer(targetContainerName)
	if err != nil {
		// retry image build if the error is image found error.
		logger.Println(err.Error(), "Retry...")
		imageExists = false
		goto BuildImage
	}

	r.lru[functionName] = time.Now().Unix()
	go r.manageCaches()
	return out, meta
}

func (r *FunctionRunner) manageCaches() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ctx := context.Background()

	live, dead := leastRecentlyUsed(&r.lru, r.cachingOptions.imageCacheMaxNumber)

	for _, functionName := range dead {
		// Remove docker image
		if r.imageExists[functionName] {
			if _, err := cli.ImageRemove(ctx, imageTagName(functionName), types.ImageRemoveOptions{Force: true}); err != nil {
				logger.Println(err.Error())
			}
		}
		delete(r.imageExists, functionName)

		// Remove pre-warm containers
		for i := len(r.prewarmedContainers)-1; i >= 0; i-- {
			containerName := r.prewarmedContainers[i]
			if containerBelongsToFunction(containerName, functionName) {
				r.prewarmedContainers = append(r.prewarmedContainers[:i], r.prewarmedContainers[i+1:]...)
				cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{})
			}
		}
	}

	for _, functionName := range live {
		r.imageExists[functionName] = true

		// Make pool
		prewarmedContainerNum := 0
		for _, containerName := range r.prewarmedContainers {
			if containerBelongsToFunction(containerName, functionName) {
				prewarmedContainerNum++
			}
		}

		// TODO
		poolSize := 4
		for i := 0; i < poolSize - prewarmedContainerNum; i++ {
			containerName, _ := CreateContainer(functionName)
			r.prewarmedContainers = append(r.prewarmedContainers, containerName)
		}
	}
}

