package scheduler

type Node struct {
	Id          int
	Url         string
	MaxCapacity int
}

func NewNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.Id = id
	n.Url = url
	n.MaxCapacity = maxCapacity
	return n
}

type Scheduler interface {
	Select(string) (*Node, error)
	Finished(*Node, string, int64) error
}
