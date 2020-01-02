package main

import "errors"

type Scheduler interface {
	pick(string) (*Node, error)
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
