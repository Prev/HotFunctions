package main

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

type FunctionRunner struct {
	imageBuilder         *ImageBuilder
	cachingOptions       CachingOptions
	lru                  map[string]int64
	imageExists          map[string]bool
	poolManager          *ContainerPoolManager
	mutex                *sync.Mutex
}

type runningMetaData struct {
	ImageBuilt         bool
	UsePooledContainer bool
}

func newFunctionRunner(cachingOptions CachingOptions) *FunctionRunner {
	r := new(FunctionRunner)
	r.imageBuilder = newImageBuilder()
	r.cachingOptions = cachingOptions
	r.lru = make(map[string]int64)
	r.imageExists = make(map[string]bool)
	r.poolManager = NewContainerPoolManager(&cachingOptions)
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
		if err = r.imageBuilder.BuildSafe(functionName); err != nil {
			logger.Println(err.Error())
		}
		meta.ImageBuilt = true
		r.imageExists[functionName] = true
	}

	// Step2: Check for the container pool.
	targetContainerName, _ := r.poolManager.Pop(functionName)

	if targetContainerName != "" {
		meta.UsePooledContainer = true
	} else {
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

	for i, functionName := range sortMapByValue(&r.lru) {
		// Check for the image cache limit
		// If the limit is set -1, the is no limitation
		if r.cachingOptions.ImageLimit != -1 && i >= r.cachingOptions.ImageLimit {
			if r.imageExists[functionName] {
				if _, err := cli.ImageRemove(ctx, imageTagName(functionName), types.ImageRemoveOptions{Force: true}); err != nil {
					logger.Println(err.Error())
				}
				delete(r.imageExists, functionName)
			}
		}

		// Check for the container pooling
		// If the limit is set -1, do not use the feature
		if r.cachingOptions.ContainerPoolLimit != -1 {
			if i < r.cachingOptions.ContainerPoolLimit {
				r.poolManager.MakePool(functionName)

			} else {
				r.poolManager.Clear(functionName)
			}
		}
	}
}

