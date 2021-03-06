package scheduler

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type ConsistentHashingScheduler struct {
	Scheduler
	virtualNodes     []vNode
	maxLoadThreshold uint
	mutex            *sync.Mutex
}

type vNode struct {
	hashkey int
	node    *Node
}

func NewConsistentHashingScheduler(nodes *[]*Node, numVirtualNodes int, maxLoadThreshold uint) *ConsistentHashingScheduler {
	s := ConsistentHashingScheduler{}
	s.maxLoadThreshold = maxLoadThreshold
	s.virtualNodes = make([]vNode, len(*nodes)*numVirtualNodes)
	s.mutex = new(sync.Mutex)

	for i, node := range *nodes {
		for m := 0; m < numVirtualNodes; m++ {
			key := fmt.Sprintf("%d-%d", node.Id, m)
			s.virtualNodes[i*numVirtualNodes+m] = vNode{s.hash(key), node}
		}
	}
	sort.SliceStable(s.virtualNodes, func(i, j int) bool {
		return s.virtualNodes[i].hashkey < s.virtualNodes[j].hashkey
	})
	return &s
}

func (s ConsistentHashingScheduler) hash(key string) int {
	b := md5.Sum([]byte(key))
	b2 := b[0:4]
	return int(binary.BigEndian.Uint32(b2)) % 1234567 // magic number
}

func (s ConsistentHashingScheduler) Select(functionName string) (*Node, error) {
	hashkey := s.hash(functionName)
	n := len(s.virtualNodes)

	// Binary search
	left := 0
	right := n
	var mid int

	for left < right {
		mid = left + int((right-left)/2)
		if s.virtualNodes[mid].hashkey < hashkey {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if right == n {
		right = 0
	}

	//return s.virtualNodes[right].node, nil

	// Bounded load
	s.mutex.Lock()
	defer s.mutex.Unlock()
	i := right
	last := (i - 1 + n) % n
	for {
		if s.virtualNodes[i].node.Load < s.maxLoadThreshold {

			s.virtualNodes[i].node.Load++
			return s.virtualNodes[i].node, nil
		}
		if i == last {
			// all ring elements are visited
			break
		}
		i = (i + 1) % n
	}

	return nil, errors.New("no available node found")
}

func (s ConsistentHashingScheduler) Finished(node *Node, _ string) error {
	s.mutex.Lock()
	node.Load--
	s.mutex.Unlock()
	return nil
}
