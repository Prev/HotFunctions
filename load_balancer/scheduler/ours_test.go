package scheduler

import (
	"testing"
)

func Test_sliceTopN(t *testing.T) {
	data := make(map[string]int)
	data["A"] = 4
	data["B"] = 5
	data["C"] = 3
	data["D"] = 1
	data["E"] = 2

	ret := sliceTopN(data, 2)
	//print(ret[0], ret[1])

	t.Log(ret[0], ret[1])
}

func Test_Ours(t *testing.T) {
	nodeList := make([]*Node, 2)
	for i := 0; i < 2; i++ {
		nodeList[i] = NewNode(i, "")
	}

	sched := NewOurScheduler(&nodeList, 4, 3, 1)

	functionStream := []string{"A", "A", "B", "C", "B", "C"}
	expectedNodes := []int{0, 0, 1, 1, 1, 0}

	for i, fn := range functionStream {
		node, err := sched.Select(fn)
		if err != nil {
			t.Fatal(err)
		}

		if node.Id != expectedNodes[i] {
			t.Fatal("Expected nodeId:", expectedNodes[i], ", real: ", node.Id)
		}
	}

	if len(sched.assigned["C"]) != 1 || sched.assigned["C"][0].Id != 0 {
		t.Fatal("The only assigned node of C should be Node0")
	}
}

// func Test_Adaptive_Finished(t *testing.T) {
// 	nodeList := make([]*Node, 2)
// 	for i := 0; i < 2; i++ {
// 		nodeList[i] = new(Node)
// 		nodeList[i].Id = i
// 		nodeList[i].MaxCapacity = 6
// 	}

// 	sched := NewAdaptiveScheduler(&nodeList, 2)

// 	n1, _ := sched.Select("A")
// 	n2, _ := sched.Select("B")
// 	sched.Finished(n1, "A", 1000)

// 	if _, ok := sched.runningTable[n1.Id]["A"]; ok == true {
// 		t.Fatal("function A should not be in the table, but it is")
// 	}

// 	if n, _ := sched.runningTable[n1.Id]["B"]; n == 1 {
// 		t.Fatal("function B should be 1, not", n)
// 	}

// 	n3, _ := sched.Select("C")
// 	sched.Finished(n3, "C", 2000)
// 	sched.Finished(n2, "B", 3000)

// 	for _, node := range nodeList {
// 		if sched.nodeCapacity(node) != 6 {
// 			t.Fatal("All function expected to be finished, but there are remaining table entries")
// 		}
// 	}
// }
