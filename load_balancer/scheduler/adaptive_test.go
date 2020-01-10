package scheduler

import (
	"testing"
)

func Test_Adaptive(t *testing.T) {
	nodeList := make([]*Node, 2)
	for i := 0; i < 2; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
		nodeList[i].MaxCapacity = 6
	}

	sched := NewAdaptiveScheduler(&nodeList, 2)

	function_stream := []string{"A", "A", "A", "B", "C", "B", "B", "C"}
	expected_nodes := []int{0, 0, 0, 1, 1, 1, 1, 0}

	for i, fn := range function_stream {
		node, err := sched.Select(fn)
		if err != nil {
			t.Fatal(err)
		}

		if node.Id != expected_nodes[i] {
			t.Fatal("Expected nodeId:", expected_nodes[i], ", real: ", node.Id)
		}
	}
}

func Test_Adaptive_Finished(t *testing.T) {
	nodeList := make([]*Node, 2)
	for i := 0; i < 2; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
		nodeList[i].MaxCapacity = 6
	}

	sched := NewAdaptiveScheduler(&nodeList, 2)

	n1, _ := sched.Select("A")
	n2, _ := sched.Select("B")
	sched.Finished(n1, "A", 1000)

	if _, ok := sched.runningTable[n1.Id]["A"]; ok == true {
		t.Fatal("function A should not be in the table, but it is")
	}

	if n, _ := sched.runningTable[n1.Id]["B"]; n == 1 {
		t.Fatal("function B should be 1, not", n)
	}

	n3, _ := sched.Select("C")
	sched.Finished(n3, "C", 2000)
	sched.Finished(n2, "B", 3000)

	for _, node := range nodeList {
		if sched.nodeCapacity(node) != 6 {
			t.Fatal("All function expected to be finished, but there are remaining table entries")
		}
	}
}
