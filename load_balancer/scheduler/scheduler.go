package scheduler

type Scheduler interface {
	Select(string) (*Node, error)
	Finished(*Node, string) error
}
