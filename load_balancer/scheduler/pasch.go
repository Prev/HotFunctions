package scheduler

import (
	"fmt"
	"sort"
	"sync"
)

const LOAD_THREASHOLD = 8

type PASchExtendedScheduler struct {
	ConsistentHashingScheduler
	connections map[int]int
	mutex       *sync.Mutex
}


func NewPASchScheduler(nodes *[]*Node) *PASchExtendedScheduler {
	s := PASchExtendedScheduler{}

	// For consistent hashing
	numVirtualNodes := 8
	s.virtualNodes = make([]vNode, len(*nodes)*numVirtualNodes)

	for i, node := range *nodes {
		for m := 0; m < numVirtualNodes; m++ {
			key := fmt.Sprintf("%d-%d", node.Id, m)
			s.virtualNodes[i*numVirtualNodes+m] = vNode{s.hash(key), node}
		}
	}
	sort.SliceStable(s.virtualNodes, func(i, j int) bool {
		return s.virtualNodes[i].hashkey < s.virtualNodes[j].hashkey
	})

	// For least loaded algorithm
	s.connections = make(map[int]int)
	for _, node := range *nodes {
		s.connections[node.Id] = 0
	}
	s.mutex = new(sync.Mutex)
	return &s
}

func (s PASchExtendedScheduler) Select(functionName string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	a1, _ := s.ConsistentHashingScheduler.Select(functionName)
	a2, _ := s.ConsistentHashingScheduler.Select(functionName + "salt")

	load1 := s.connections[a1.Id]
	load2 := s.connections[a2.Id]

	var selected *Node = nil

	if load1 < load2 {
		selected = a1
	} else {
		selected = a2
	}

	if load1 > LOAD_THREASHOLD && load2 > LOAD_THREASHOLD {
		minLoaded := 9999999
		for _, vn := range s.virtualNodes {
			loaded := s.connections[vn.node.Id]
			if minLoaded > loaded {
				minLoaded = loaded
				selected = vn.node
			}
		}
	}

	if _, exists := s.connections[selected.Id]; exists == false {
		s.connections[selected.Id] = 0
	}
	s.connections[selected.Id] += 1
	return selected, nil
}

func (s PASchExtendedScheduler) Finished(node *Node, _ string, _ int64) error {
	s.mutex.Lock()
	s.connections[node.Id] -= 1
	s.mutex.Unlock()

	return nil
}