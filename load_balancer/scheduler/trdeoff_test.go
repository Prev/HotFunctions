package scheduler

import (
	"testing"
)

func Test_TradeOff(t *testing.T) {
	nodeList := make([]*Node, 10)
	for i := 0; i < 10; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
	}

	// Acts like least loaded sched
	llSched := NewTradeOffScheduler(&nodeList, 0, 0)

	for i := 0; i < 80; i++ {
		node, err := llSched.Select("")
		if err != nil {
			t.Fatal(err)
		}
		if node.Id != i%10 {
			t.Fatal("Load balanced not by the least loaded algorithm")
		}
	}

	// Acts like hash sched
	hashSched := NewTradeOffScheduler(&nodeList, 0, 1)

	var a1, a2, a3, b1, b2 *Node
	var err error
	if a1, err = hashSched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if b1, err = hashSched.Select("B"); err != nil {
		t.Fatal(err)
	}
	if _, err = hashSched.Select("C"); err != nil {
		t.Fatal(err)
	}
	if a2, err = hashSched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if a3, err = hashSched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if _, err = hashSched.Select("A"); err != nil {
		t.Fatal(err)
	}
	if b2, err = hashSched.Select("B"); err != nil {
		t.Fatal(err)
	}

	if a1.Id != a2.Id || a2.Id != a3.Id || b1.Id != b2.Id {
		t.Fatal("Same functions are placed different nodes")
	}
}
