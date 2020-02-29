package scheduler

import (
	"testing"
)

func Test_RoundRobin(t *testing.T) {
	nodeList := make([]*Node, 10)
	for i := 0; i < 10; i++ {
		nodeList[i] = new(Node)
		nodeList[i].Id = i
	}

	sched := NewRoundRobinScheduler(&nodeList)

	for i := 0; i < 80; i++ {
		node, err := sched.Select("")
		if err != nil {
			t.Fatal(err)
		}
		if node.Id != i%10 {
			t.Fatal("Load balanced not by the round robin algorithm")
		}
	}
}
