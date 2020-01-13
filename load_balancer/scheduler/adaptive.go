package scheduler

import (
	"errors"
	"fmt"
	"sync"
)

// Greedy scheduler performing reassignment

type AdaptiveScheduler struct {
	Scheduler
	nodes             *[]*Node
	runningTable      map[int]map[string]int
	assignedTable     map[string][]int
	sumExecutionTimes map[string]int64
	numExecutions     map[string]int64
	CMin              int
	mutex             *sync.Mutex
}

func NewAdaptiveScheduler(nodes *[]*Node, CMin int) *AdaptiveScheduler {
	s := AdaptiveScheduler{}
	s.assignedTable = make(map[string][]int)
	s.sumExecutionTimes = make(map[string]int64)
	s.numExecutions = make(map[string]int64)
	s.nodes = nodes
	s.CMin = CMin
	s.mutex = new(sync.Mutex)

	s.runningTable = make(map[int]map[string]int)
	for _, node := range *s.nodes {
		s.runningTable[node.Id] = make(map[string]int)
	}

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

func (s AdaptiveScheduler) Finished(node *Node, functionName string, executionTime int64) error {
	s.mutex.Lock()
	s.runningTable[node.Id][functionName]--
	if s.runningTable[node.Id][functionName] == 0 {
		delete(s.runningTable[node.Id], functionName)
	}

	s.sumExecutionTimes[functionName] += executionTime
	s.numExecutions[functionName]++
	s.mutex.Unlock()
	return nil
}

func (s AdaptiveScheduler) Select(functionName string) (*Node, error) {
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
				if node.Id == nodeId {
					candidates = append(candidates, node)
					break
				}
			}
		}
		selected, _ = s.maxCapacityNode(&candidates)
	}

	if selected == nil {
		selected, err = s.maxCapacityNode(s.nodes)
		if err != nil {
			return nil, err
		}

		// Register for future use
		s.assignedTable[functionName] = append(s.assignedTable[functionName], selected.Id)
	}

	if _, exists := s.runningTable[selected.Id][functionName]; exists == false {
		s.runningTable[selected.Id][functionName] = 0
	}
	s.runningTable[selected.Id][functionName] += 1

	s.performReassignment(selected)
	return selected, nil
}

func (s AdaptiveScheduler) nodeCapacity(node *Node) int {
	running := 0
	for _, val := range s.runningTable[node.Id] {
		running += val
	}
	return node.MaxCapacity - running
}

func (s AdaptiveScheduler) maxCapacityNode(candidates *[]*Node) (*Node, error) {
	var selected *Node = nil
	maxCapacity := 0

	for _, node := range *candidates {
		capacity := s.nodeCapacity(node)

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

func (s AdaptiveScheduler) performReassignment(node *Node) {
	capacity := s.nodeCapacity(node)
	if capacity > s.CMin {
		// does not violate the policy
		return
	}

	majorityFname := ""
	majorityVal := 0

	for fname, val := range s.runningTable[node.Id] {
		if majorityVal < val {
			majorityVal = val
			majorityFname = fname
		}
	}

	// pick function with shortest MET except for the majority function
	shortestMet := 999999999
	victimFunction := ""

	for fname, _ := range s.runningTable[node.Id] {
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
			if nodeId == node.Id {
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
