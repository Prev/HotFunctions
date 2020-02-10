package main

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types"
)

type ContainerPoolManager struct {
	poolSize int
	pool     []string
	mutex    *sync.Mutex
}

func NewContainerPoolManager(poolSize int) *ContainerPoolManager {
	m := new(ContainerPoolManager)
	m.poolSize = poolSize
	m.pool = make([]string, 0)
	m.mutex = new(sync.Mutex)
	return m
}

// Pick one of the container pool of the function.
// If pool not exists, return "", false
func (m *ContainerPoolManager) Pop(functionName string) (string, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, containerName := range m.pool {
		if containerBelongsToFunction(containerName, functionName) {
			// Pick the proper container and remove it from the pool
			m.pool = append(m.pool[:i], m.pool[i+1:]...)
			return containerName, true
		}
	}
	return "", false
}

// Make the container pool of the function
func (m *ContainerPoolManager) MakePool(functionName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	existentPoolNum := 0
	for _, containerName := range m.pool {
		if containerBelongsToFunction(containerName, functionName) {
			existentPoolNum++
		}
	}

	for i := 0; i < m.poolSize - existentPoolNum; i++ {
		containerName, _ := CreateContainer(functionName)
		m.pool = append(m.pool, containerName)
	}
}

// Clear the container pool of the function
func (m *ContainerPoolManager) Clear(functionName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ctx := context.Background()

	// Iterate reversed order to make the multiple deletion work properly
	for i := len(m.pool)-1; i >= 0; i-- {
		containerName := m.pool[i]
		if containerBelongsToFunction(containerName, functionName) {
			m.pool = append(m.pool[:i], m.pool[i+1:]...)
			cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{})
		}
	}
}
