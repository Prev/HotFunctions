package scheduler

import (
	"errors"
	"sync"
)

type LeastLoadedScheduler struct {
	Scheduler
	nodes       *[]*Node
	mutex       *sync.Mutex
}

func NewLeastLoadedScheduler(nodes *[]*Node) *LeastLoadedScheduler {
	s := LeastLoadedScheduler{}
	s.nodes = nodes
	s.mutex = new(sync.Mutex)
	return &s
}

func (s LeastLoadedScheduler) Select(_ string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var selected *Node = nil
	for _, node := range *s.nodes {
		if selected == nil || node.Load < selected.Load {
			selected = node
		}
	}

	if selected == nil {
		return nil, errors.New("no available node found")
	}
	selected.Load++
	return selected, nil
}

func (s LeastLoadedScheduler) Finished(node *Node, _ string) error {
	s.mutex.Lock()
	node.Load--
	s.mutex.Unlock()
	return nil
}
