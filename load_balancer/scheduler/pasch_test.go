package scheduler

import "testing"

func Test_PASchExtended(t *testing.T) {
	nodeList := make([]*Node, 6)
	for i := 0; i < 6; i++ {
		nodeList[i] = NewNode(i, "")
	}

	sched := NewPASchScheduler(&nodeList, 8)

	functionStream := []string{"A", "A", "A", "B", "C", "B", "B", "C"}
	expectedNodes := []int{5, 3, 5, 0, 1, 0, 5, 1}

	for i, fn := range functionStream {
		node, err := sched.Select(fn)
		if err != nil {
			t.Fatal(err)
		}

		//t.Log(i, fn, node.Id)
		if node.Id != expectedNodes[i] {
			t.Fatal("Expected nodeId:", expectedNodes[i], ", real: ", node.Id)
		}
	}
}