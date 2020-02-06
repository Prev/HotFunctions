package scheduler

import (
	"strconv"
	"sync"

	"github.com/lafikl/consistent"
)

type PASchExtendedScheduler struct {
	Scheduler
	hashRing      *consistent.Consistent
	loadThreshold uint
	workerNodes   *[]*Node
	workerNodeMap map[string]*Node
	mutex         *sync.Mutex
}

func NewPASchExtendedScheduler(nodes *[]*Node, loadThreshold uint) *PASchExtendedScheduler {
	hashRing := consistent.New()
	workerNodeMap := make(map[string]*Node)

	s := PASchExtendedScheduler{
		nil,
		hashRing,
		loadThreshold,
		nodes,
		workerNodeMap,
		new(sync.Mutex),
	}

	for _, node := range *nodes {
		key := strconv.Itoa(node.Id)
		s.hashRing.Add(key)
		s.workerNodeMap[key] = node
	}
	return &s
}

func (s *PASchExtendedScheduler) Select(functionName string) (*Node, error) {
	key1, _ := s.hashRing.Get(functionName)
	key2, _ := s.hashRing.Get(functionName + "salt")

	node1 := s.workerNodeMap[key1]
	node2 := s.workerNodeMap[key2]

	selectedNode := node1
	if node1.Load > node2.Load {
		selectedNode = node2
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if selectedNode.Load >= s.loadThreshold { // Find least loaded
		selectedNode = s.selectLeastLoadedWorker()
	}

	selectedNode.Load++
	return selectedNode, nil
}

func (s *PASchExtendedScheduler) Finished(node *Node, _ string) error {
	s.mutex.Lock()
	node.Load--
	s.mutex.Unlock()
	return nil
}

func (s *PASchExtendedScheduler) selectLeastLoadedWorker() *Node {
	var selected *Node = nil
	for _, node := range *s.workerNodes {
		if selected == nil || node.Load < selected.Load {
			selected = node
		}
	}
	return selected
}
