package scheduler

import (
	"errors"
	"sync"
)

type LeastLoadedScheduler struct {
	Scheduler
	nodes         *[]*Node
	capacityTable map[int]int
	mutex         *sync.Mutex
}

func NewLeastLoadedScheduler(nodes *[]*Node) *LeastLoadedScheduler {
	s := LeastLoadedScheduler{}
	s.nodes = nodes
	s.mutex = new(sync.Mutex)

	// Init capacity table
	s.capacityTable = make(map[int]int)
	for _, node := range *s.nodes {
		s.capacityTable[node.Id] = node.MaxCapacity
	}
	return &s
}

func (s LeastLoadedScheduler) Select(_ string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var selected *Node = nil
	minUsed := 999999

	for _, node := range *s.nodes {
		used := node.MaxCapacity - s.capacityTable[node.Id]
		if used >= node.MaxCapacity {
			continue
		}
		if used < minUsed {
			minUsed = used
			selected = node
		}
	}

	if selected == nil {
		return nil, errors.New("no available node found")
	}

	s.capacityTable[selected.Id] -= 1
	return selected, nil
}

func (s LeastLoadedScheduler) Finished(node *Node, _ string, _ int64) error {
	s.mutex.Lock()
	s.capacityTable[node.Id] += 1
	s.mutex.Unlock()

	return nil
}
