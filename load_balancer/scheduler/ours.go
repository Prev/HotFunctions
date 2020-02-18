package scheduler

import (
	"errors"
	"sort"
	"sync"
)

type OurScheduler struct {
	Scheduler
	nodes     *[]*Node
	assigned  map[string][]*Node
	running   map[int]map[string]int
	TMax      uint
	TOpt      uint
	CacheSize int
	mutex     *sync.Mutex
}

func NewOurScheduler(nodes *[]*Node, TMax uint, TOpt uint, cacheSize int) *OurScheduler {
	s := OurScheduler{}
	s.nodes = nodes
	s.TMax = TMax
	s.TOpt = TOpt
	s.CacheSize = cacheSize
	s.assigned = make(map[string][]*Node)
	s.mutex = new(sync.Mutex)

	s.running = make(map[int]map[string]int)
	for _, node := range *s.nodes {
		s.running[node.Id] = make(map[string]int)
	}
	return &s
}

func (s *OurScheduler) Select(functionName string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var selected *Node = nil

	candidateNodes, exists := s.assigned[functionName]
	if exists == false {
		s.assigned[functionName] = make([]*Node, 0)
	}

	if len(candidateNodes) > 0 {
		selected = s.leastLoadedAmongAvailable(functionName, &candidateNodes)
	}

	if selected == nil {
		if selected = s.leastLoadedAmongAvailable(functionName, s.nodes); selected == nil {
			return nil, errors.New("no available node found")
		}
		// Register for future use
		s.assigned[functionName] = append(s.assigned[functionName], selected)
	}

	if _, exists := s.running[selected.Id][functionName]; exists == false {
		s.running[selected.Id][functionName] = 0
	}
	s.running[selected.Id][functionName] += 1
	selected.Load++
	return selected, nil
}

func (s *OurScheduler) Finished(node *Node, functionName string) error {
	s.mutex.Lock()
	s.running[node.Id][functionName]--
	node.Load--
	s.mutex.Unlock()
	return nil
}

func (s *OurScheduler) available(node *Node, f string) bool {
	if node.Load >= s.TMax {
		// Task is overloaded
		return false

	} else if node.Load >= s.TOpt {
		// Work load is going full, only accept for the major applications
		majorFunctions := sliceTopN(s.running[node.Id], s.CacheSize)
		for _, fi := range majorFunctions {
			if fi == f {
				return true
			}
		}
		// De-assign node from assigned table
		idx := -1
		for i, n := range s.assigned[f] {
			if n == node {
				idx = i
			}
		}
		if idx != -1 {
			s.assigned[f] = append(s.assigned[f][:idx], s.assigned[f][idx+1:]...)
		}
		return false

	} else {
		return true
	}
}

func (s *OurScheduler) leastLoadedAmongAvailable(functionName string, candidates *[]*Node) *Node {
	var selected *Node = nil
	for _, node := range *candidates {
		if !s.available(node, functionName) {
			continue
		}
		if selected == nil || node.Load < selected.Load {
			selected = node
		}
	}
	return selected
}

func sliceTopN(data map[string]int, n int) []string {
	values := make([]int, 0)
	for _, ni := range data {
		values = append(values, ni)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})

	ret := make([]string, 0)
	for i, ni := range values {
		if i >= n {
			break
		}
		for fj, nj := range data {
			if ni == nj {
				ret = append(ret, fj)
				break
			}
		}
	}
	return ret
}
