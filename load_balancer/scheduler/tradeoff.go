package scheduler

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math/rand"
	"sync"
)

type TradeOffScheduler struct {
	Scheduler
	nodes *[]*Node
	mutex *sync.Mutex
	alpha float64
	beta  float64
}

func NewTradeOffScheduler(nodes *[]*Node, alpha float64, beta float64) *TradeOffScheduler {
	s := TradeOffScheduler{}
	s.nodes = nodes
	s.mutex = new(sync.Mutex)
	s.alpha = alpha
	s.beta = beta
	return &s
}

func (s *TradeOffScheduler) hash(key string) int {
	b := md5.Sum([]byte(key))
	b2 := b[0:4]
	return int(binary.BigEndian.Uint32(b2)) % 1234567 // magic number
}

func (s *TradeOffScheduler) Select(functionName string) (*Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var selected *Node = nil

	if rand.Float64() > 1-s.alpha {
		// Random selection
		selected = (*s.nodes)[rand.Intn(len(*s.nodes))]

	} else if rand.Float64() > 1-s.beta {
		// Hash
		selected = (*s.nodes)[s.hash(functionName)%len(*s.nodes)]

	} else {
		// Least Loaded
		for _, node := range *s.nodes {
			if selected == nil || node.Load < selected.Load {
				selected = node
			}
		}
	}

	if selected == nil {
		return nil, errors.New("no available node found")
	}
	selected.Load++
	return selected, nil
}

func (s *TradeOffScheduler) Finished(node *Node, _ string) error {
	s.mutex.Lock()
	node.Load--
	s.mutex.Unlock()
	return nil
}
