package main

import "errors"

func leastLoaded(nodes *[]*Node) (*Node, error) {
	selected := -1
	minUsed := 999999

	for i, node := range *nodes {
		used := node.maxCapacity - node.capacity()
		if used >= node.maxCapacity {
			continue
		}
		if used < minUsed {
			minUsed = used
			selected = i
		}
	}

	if selected == -1 {
		return nil, errors.New("no available node found")
	}

	return (*nodes)[selected], nil
}

func maxCapacity(nodes *[]*Node) (*Node, error) {
	selected := -1
	maxCapacity := 0

	for i, node := range *nodes {
		capacity := node.capacity()

		if capacity <= 0 {
			continue
		}
		if capacity > maxCapacity {
			maxCapacity = capacity
			selected = i
		}
	}

	if selected == -1 {
		return nil, errors.New("no available node found")
	}

	return (*nodes)[selected], nil
}
