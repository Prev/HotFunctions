package main

import "errors"

type LeastLoadedScheduler struct {
	nodes *[]*Node
}

func newLeastLoadedScheduler(nodes *[]*Node) *LeastLoadedScheduler {
	s := LeastLoadedScheduler{}
	s.nodes = nodes
	return &s
}

func (s LeastLoadedScheduler) pick(_ string) (*Node, error) {
	selected := -1
	minUsed := 999999

	for i, node := range *s.nodes {
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

	return (*s.nodes)[selected], nil
}
