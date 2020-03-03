package scheduler

import (
	"testing"
)

func Test_ConsistentHashing(t *testing.T) {
	nodeList := make([]*Node, 4)
	for i := 0; i < 4; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
	}

	sched := NewConsistentHashingScheduler(&nodeList, 4, 4)

	var n1, n2, n3, n4 *Node
	var err error

	if n1, err = sched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if _, err = sched.Select("B"); err != nil {
		t.Fatal(err)
	}
	if _, err = sched.Select("C"); err != nil {
		t.Fatal(err)
	}
	if n2, err = sched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if n3, err = sched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if _, err = sched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if n4, err = sched.Select("A"); err != nil {
		t.Fatal(err)
	}

	if n1.Id != n2.Id || n1.Id != n3.Id {
		t.Fatal("Same functions are placed different nodes")
	}

	if n4.Id == n3.Id {
		t.Fatal("Something goes wrong on processing bounded load")
	}
}
