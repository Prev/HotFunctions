package scheduler

import (
	"testing"
)

func Test_LeastLoaded(t *testing.T) {
	nodeList := make([]*Node, 10)
	for i := 0; i < 10; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
		//nodeList[i].MaxCapacity = 8
	}

	sched := NewLeastLoadedScheduler(&nodeList)

	for i := 0; i < 80; i++ {
		node, err := sched.Select("")
		if err != nil {
			t.Fatal(err)
		}
		if node.Id != i%10 {
			t.Fatal("Load balanced not by the least loaded algorithm")
		}
	}
}
//
//func Test_LeastLoaded_Heterogeneous(t *testing.T) {
//	nodeList := make([]*Node, 10)
//	for i := 0; i < 10; i++ {
//		n := new(Node)
//		n.Id = i
//		if i < 5 {
//			n.MaxCapacity = 4
//		} else {
//			n.MaxCapacity = 8
//		}
//		nodeList[i] = n
//	}
//
//	sched := NewLeastLoadedScheduler(&nodeList)
//
//	for i := 0; i < 40; i++ {
//		node, err := sched.Select("")
//		if err != nil {
//			t.Fatal(err)
//		}
//		if node.Id != i%10 {
//			t.Fatal("Load balanced not by the least loaded algorithm")
//		}
//	}
//	for i := 0; i < 20; i++ {
//		node, err := sched.Select("")
//		if err != nil {
//			t.Fatal(err)
//		}
//		if node.Id != i%5+5 {
//			t.Fatal("Load balanced not by the least loaded algorithm")
//		}
//	}
//}
//
//func Test_LeastLoaded_Finished(t *testing.T) {
//	nodeList := make([]*Node, 10)
//	for i := 0; i < 10; i++ {
//		nodeList[i] = new(Node)
//		nodeList[i].Id = i
//		nodeList[i].MaxCapacity = 8
//	}
//
//	sched := NewLeastLoadedScheduler(&nodeList)
//
//	n1, _ := sched.Select("A")
//	n2, _ := sched.Select("B")
//	sched.Finished(n1, "A", 0)
//
//	if sched.connections[n1.Id] != 8 {
//		t.Fatal("Least loaded algorithm has an error")
//	}
//
//	if sched.connections[n2.Id] != 7 {
//		t.Fatal("Least loaded algorithm has an error")
//	}
//
//	n3, _ := sched.Select("C")
//	sched.Finished(n3, "C", 0)
//	sched.Finished(n2, "B", 0)
//
//	for _, node := range nodeList {
//		if sched.connections[node.Id] != 8 {
//			t.Fatal("Least loaded algorithm has an error")
//		}
//	}
//}
