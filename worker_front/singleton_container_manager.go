package main

import (
	"sync"
	"time"
)

type SingletonContainerManager struct {
	isCreating  map[string]bool
	containers  map[string]Container
	mutex       *sync.Mutex
}

func newRestContainerManager() *SingletonContainerManager {
	b := new(SingletonContainerManager)
	b.isCreating = make(map[string]bool)
	b.containers = make(map[string]Container)
	b.mutex = new(sync.Mutex)
	return b
}

func (b *SingletonContainerManager) Get(image Image) (Container, bool) {
	b.mutex.Lock()
	cont, exists := b.containers[image.FunctionName]
	b.mutex.Unlock()
	return cont, exists
}

func (b *SingletonContainerManager) SafeCreate(image Image) (Container, error) {
	functionName := image.FunctionName
	b.mutex.Lock()

	if cont, exists := b.containers[functionName]; exists {
		b.mutex.Unlock()
		return cont, nil

	} else if b.isCreating[functionName] == true {
		// Wait until container is created
		for {
			b.mutex.Unlock()
			time.Sleep(time.Second / 20)
			b.mutex.Lock()

			if b.isCreating[functionName] == false {
				b.mutex.Unlock()
				return b.containers[functionName], nil
			}
		}

	} else {
		b.isCreating[functionName] = true
		b.mutex.Unlock()

		cont, err := CreateContainer(image)
		if err != nil {
			return Container{}, err
		}

		b.mutex.Lock()
		b.isCreating[functionName] = false
		b.containers[functionName] = cont
		b.mutex.Unlock()

		return cont, nil
	}
}

func (b *SingletonContainerManager) Delete(functionName string) {
	b.mutex.Lock()
	delete(b.containers, functionName)
	b.mutex.Unlock()
}
