package main

import "fmt"

type TableBasedScheduler struct {
	lookupTable map[string][]int // function-nodes pair
}

func newTableBasedScheduler() *TableBasedScheduler {
	s := TableBasedScheduler{}
	s.lookupTable = make(map[string][]int)
	return &s
}

func (s TableBasedScheduler) pick(nodes *[]*Node, functionName string) (*Node, error) {
	var selected *Node = nil
	var err error

	nodeIdList, exists := s.lookupTable[functionName]
	if exists == false {
		s.lookupTable[functionName] = make([]int, 0)
	}

	if exists == true {
		candidates := make([]*Node, 0, len(*nodes))
		for _, nodeId := range nodeIdList {
			for _, node := range *nodes {
				if node.id == nodeId {
					candidates = append(candidates, node)
					break
				}
			}
		}
		// Select from candidates
		selected, _ = maxCapacity(&candidates)
	}

	if selected == nil {
		selected, err = maxCapacity(nodes)
		if err != nil {
			return nil, err
		}

		// Register for future use
		s.lookupTable[functionName] = append(s.lookupTable[functionName], selected.id)
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
