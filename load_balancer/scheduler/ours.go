package scheduler

import (
	"math/rand"
	"sync"
	"time"
)

// Greedy scheduler

type OurScheduler struct {
	Scheduler
	nodes *[]*Node
	T     map[string][]*Node
	CMin  int
	mutex *sync.Mutex
}

func NewOurScheduler(nodes *[]*Node, CMin int) *OurScheduler {
	s := OurScheduler{}
	s.T = make(map[string][]*Node)
	s.nodes = nodes
	s.CMin = CMin
	s.mutex = new(sync.Mutex)
	rand.Seed(time.Now().Unix())
	return &s
}

func (s OurScheduler) Finished(node *Node, functionName string, executionTime int64) error {
	s.mutex.Lock()
	node.running[functionName] -= 1
	s.mutex.Unlock()
	return nil
}

func (s OurScheduler) Select(functionName string) (*Node, error) {
	var selected *Node = nil

	s.mutex.Lock()
	defer s.mutex.Unlock()

	delegatedNodes, exists := s.T[functionName]
	if exists == false {
		s.T[functionName] = make([]*Node, 0)
	}

	if exists == true {
		w1 := randSelect(delegatedNodes)
		w2 := randSelect(delegatedNodes)

		available1, load1 := w1.preflight(functionName)
		available2, load2 := w2.preflight(functionName)

		if available1 || available2 {
			if load1 < load2 {
				selected = w1
			} else {
				selected = w2
			}
		}
	}

	if selected == nil {
		w1 := randSelect(*s.nodes)
		w2 := randSelect(*s.nodes)

		_, load1 := w1.preflight(functionName)
		_, load2 := w2.preflight(functionName)

		if load1 < load2 {
			selected = w1
		} else {
			selected = w2
		}

		// Register for future use
		s.T[functionName] = append(s.T[functionName], selected)
	}

	if _, exists := selected.running[functionName]; exists == false {
		selected.running[functionName] = 0
	}
	selected.running[functionName] += 1
	return selected, nil
}

func randSelect(nodeList []*Node) *Node {
	i := rand.Int() % len(nodeList)
	return nodeList[i]
}
