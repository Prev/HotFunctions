package scheduler

import (
	"testing"
)

func Test_TableBased(t *testing.T) {
	nodeList := make([]*Node, 4)
	for i := 0; i < 4; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
		nodeList[i].MaxCapacity = 4
	}

	function_stream := []string{"A", "B", "A", "A", "A", "A", "C", "D", "A"}
	expected_nodes := []int{0, 1, 0, 0, 0, 2, 3, 1, 2}

	sched := NewTableBasedScheduler(&nodeList)

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
