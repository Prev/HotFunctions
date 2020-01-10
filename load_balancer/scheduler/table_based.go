package scheduler

import (
	"errors"
	"fmt"
	"sync"
)

type TableBasedScheduler struct {
	Scheduler
	nodes         *[]*Node
	capacityTable map[int]int
	lookupTable   map[string][]int // function-nodes
	mutex         *sync.Mutex
}

func NewTableBasedScheduler(nodes *[]*Node) *TableBasedScheduler {
	s := TableBasedScheduler{}
	s.nodes = nodes
	s.lookupTable = make(map[string][]int)
	s.mutex = new(sync.Mutex)

	// Init capacity table
	s.capacityTable = make(map[int]int)
	for _, node := range *s.nodes {
		s.capacityTable[node.Id] = node.MaxCapacity
	}
	return &s
}

func (s TableBasedScheduler) Select(functionName string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var selected *Node = nil
	var err error

	nodeIdList, exists := s.lookupTable[functionName]
	if exists == false {
		s.lookupTable[functionName] = make([]int, 0)
	}

	if exists == true {
		candidates := make([]*Node, 0, len(*s.nodes))
		for _, nodeId := range nodeIdList {
			for _, node := range *s.nodes {
				if node.Id == nodeId {
					candidates = append(candidates, node)
					break
				}
			}
		}
		// Select from candidates
		selected, _ = s.maxCapacityNode(&candidates)
	}

	if selected == nil {
		selected, err = s.maxCapacityNode(s.nodes)
		if err != nil {
			return nil, err
		}

		// Register for future use
		s.lookupTable[functionName] = append(s.lookupTable[functionName], selected.Id)
	}

	s.capacityTable[selected.Id] -= 1
	return selected, nil
}

func (s TableBasedScheduler) Finished(node *Node, _ string, _ int64) error {
	s.mutex.Lock()
	s.capacityTable[node.Id] += 1
	s.mutex.Unlock()

	return nil
}

func (s TableBasedScheduler) maxCapacityNode(candidates *[]*Node) (*Node, error) {
	var selected *Node = nil
	maxCapacity := 0

	for _, node := range *candidates {
		capacity := s.capacityTable[node.Id]

		if capacity <= 0 {
			continue
		}
		if capacity > maxCapacity {
			maxCapacity = capacity
			selected = node
		}
	}

	if selected == nil {
		return nil, errors.New("no available node found")
	}

	return selected, nil
}

func (s TableBasedScheduler) printTables() {
	out := "-------lookup table-------\n" +
		"|Func\t|Nodes\t|\n"

	for fname, nodeIdList := range s.lookupTable {
		nodeListStr := ""
		for _, id := range nodeIdList {
			nodeListStr += fmt.Sprintf("%d, ", id)
		}
		out += fmt.Sprintf("|%s\t|%s\t|\n", fname, nodeListStr)
	}
	out += "-------------------------"
	println(out)
}
