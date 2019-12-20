package main

import (
	"fmt"
	"sync"
)

// Greedy scheduler performing reassingment

type AdaptiveScheduler struct {
	assignedTable     map[string][]int
	sumExecutionTimes map[string]int64
	numExecutions     map[string]int64
	CMin              int
	mutex             *sync.Mutex
	nodes             *[]*Node
}

func newAdaptiveScheduler(nodes *[]*Node, CMin int) *AdaptiveScheduler {
	s := AdaptiveScheduler{}
	s.assignedTable = make(map[string][]int)
	s.sumExecutionTimes = make(map[string]int64)
	s.numExecutions = make(map[string]int64)
	s.nodes = nodes
	s.CMin = CMin
	s.mutex = new(sync.Mutex)
	return &s
}

func (s AdaptiveScheduler) initMapForFunction(functionName string) {
	s.assignedTable[functionName] = make([]int, 0)
	s.sumExecutionTimes[functionName] = 0
	s.numExecutions[functionName] = 0
}

func (s AdaptiveScheduler) met(functionName string) int64 {
	if s.numExecutions[functionName] == 0 {
		return 0
	}
	return s.sumExecutionTimes[functionName] / s.numExecutions[functionName]
}

func (s AdaptiveScheduler) appendExecutionResult(functionName string, executionTime int64) {
	s.mutex.Lock()
	s.sumExecutionTimes[functionName] += executionTime
	s.numExecutions[functionName]++
	s.mutex.Unlock()
}

func (s AdaptiveScheduler) pick(functionName string) (*Node, error) {
	var selected *Node = nil
	var err error

	s.mutex.Lock()
	defer s.mutex.Unlock()

	nodeIdList, exists := s.assignedTable[functionName]
	if exists == false {
		s.initMapForFunction(functionName)
	}

	if exists == true {
		candidates := make([]*Node, 0, len(*s.nodes))
		for _, nodeId := range nodeIdList {
			for _, node := range *s.nodes {
				if node.id == nodeId {
					candidates = append(candidates, node)
					break
				}
			}
		}
		selected, _ = maxCapacity(&candidates)
	}

	if selected == nil {
		selected, err = maxCapacity(s.nodes)
		if err != nil {
			return nil, err
		}

		// Register for future use
		s.assignedTable[functionName] = append(s.assignedTable[functionName], selected.id)
	}

	s.performReassignment(selected)

	return selected, nil
}

func (s AdaptiveScheduler) performReassignment(node *Node) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	capacity := node.capacity()
	if capacity > s.CMin {
		// does not violate the policy
		return
	}

	majorityFname := ""
	majorityVal := 0

	for fname, val := range node.running {
		if majorityVal < val {
			majorityVal = val
			majorityFname = fname
		}
	}

	// pick function with shortest MET except for the majority function
	shortestMet := 999999999
	victimFunction := ""

	for fname, _ := range node.running {
		if fname == majorityFname {
			// skip for the majority function in the node
			continue
		}
		met := int(s.met(fname))
		if met < shortestMet {
			shortestMet = met
			victimFunction = fname
		}
	}

	if victimFunction != "" {
		for i, nodeId := range s.assignedTable[victimFunction] {
			if nodeId == node.id {
				// delete node
				s.assignedTable[victimFunction] = append(s.assignedTable[victimFunction][:i], s.assignedTable[victimFunction][i+1:]...)
			}
		}
	}
}

func (s AdaptiveScheduler) printTables() {
	out := "-------lookup table-------\n" +
		"|Func\t|Nodes\t|MET\t|\n"

	for fname, nodeIdList := range s.assignedTable {
		nodeListStr := ""
		for _, id := range nodeIdList {
			nodeListStr += fmt.Sprintf("%d, ", id)
		}
		out += fmt.Sprintf("|%s\t|%s\t|%d\t|\n", fname, nodeListStr, s.met(fname))
	}
	out += "-------------------------"
	println(out)
}
