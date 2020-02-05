package scheduler

import "sort"

const T_MAX = 8
const T_OPT = 6
const CACHE_SIZE = 3

type Node struct {
	Id          int
	Url         string
	MaxCapacity int
	running     map[string]int
}

func NewNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.Id = id
	n.Url = url
	n.MaxCapacity = maxCapacity
	n.running = make(map[string]int)
	return n
}

func (node Node) preflight(f string) (bool, int) {
	loaded := 0
	for _, ni := range node.running {
		loaded += ni
	}

	if loaded >= T_MAX {
		// Task is overloaded
		return false, loaded

	} else if loaded >= T_OPT {
		// Work load is going full, only accept for the major applications
		majorFunctions := sliceTopN(node.running, CACHE_SIZE)
		for _, fi := range majorFunctions {
			if fi == f {
				goto ACCEPT
			}
		}
		return false, loaded

	}
	ACCEPT:
	return true, loaded
}

func sliceTopN(data map[string]int, n int) []string {
	values := make([]int, 0)
	for _, ni := range data {
		values = append(values, ni)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})

	ret := make([]string, 0)
	for i, ni := range values {
		if i >= CACHE_SIZE {
			break
		}
		for fj, nj := range data {
			if ni == nj {
				ret = append(ret, fj)
				break
			}
		}
	}
	return ret
}

type Scheduler interface {
	Select(string) (*Node, error)
	Finished(*Node, string, int64) error
}
