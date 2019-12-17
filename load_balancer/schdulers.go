package main

import "errors"

func leastLoaded(nodes *[]*Node) (*Node, error) {
	selected := -1
	minUsed := 999999

	for i, node := range *nodes {
		used := len(node.running)
		if used >= node.maxCapacity {
			continue
		}
		if used < minUsed {
			minUsed = used
			selected = i
		}
	}

	if selected == -1 {
		// panic("no available node found")
		return nil, errors.New("no available node found")
	}

	return (*nodes)[selected], nil
}
