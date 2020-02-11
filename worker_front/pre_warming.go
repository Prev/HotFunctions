package main

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Create the pre-warmed container of the image
// [Warning] make sure to call this function in the critical section
func createPreWarmedContainers(containers *[]Container, image Image, poolNum int) {
	existentContainerNum := 0
	for _, container := range *containers {
		if container.FunctionName == image.FunctionName {
			existentContainerNum++
		}
	}

	for i := 0; i < poolNum - existentContainerNum; i++ {
		container, _ := CreateContainer(image)
		*containers = append(*containers, container)
	}
}

// Clear the pre-warmed container of the image
// [Warning] make sure to call this function in the critical section
func clearPreWarmedContainers(containers *[]Container, image Image) {
	ctx := context.Background()

	// Iterate reversed order to make the multiple deletion work properly
	for i := len(*containers)-1; i >= 0; i-- {
		container := (*containers)[i]
		if container.FunctionName == image.FunctionName {
			*containers = append((*containers)[:i], (*containers)[i+1:]...)
			cli.ContainerRemove(ctx, container.Name, types.ContainerRemoveOptions{})
		}
	}
}
