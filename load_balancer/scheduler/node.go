package scheduler

type Node struct {
	Id   int
	Url  string
	Load uint
}

func NewNode(id int, url string) *Node {
	n := new(Node)
	n.Id = id
	n.Url = url
	n.Load = 0
	return n
}
