package scheduler

type RoundRobinScheduler struct {
	Scheduler
	nodes       *[]*Node
	clock       int
}

func NewRoundRobinScheduler(nodes *[]*Node) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		nodes: nodes,
		clock: 0,
	}
}

func (s* RoundRobinScheduler) Select(_ string) (*Node, error) {
	selected := (*s.nodes)[s.clock]
	s.clock = (s.clock + 1) % len(*s.nodes)
	return selected, nil
}

func (s* RoundRobinScheduler) Finished(node *Node, _ string) error {
	// do nothing
	return nil
}
