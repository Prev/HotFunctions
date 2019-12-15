package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Node struct {
	id               int
	url              string
	maxCapacity      int
	running          []string
	distinctFuncsLog []int
	warmColdLogs     []bool // true: warm, false: cold
}

func newNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.id = id
	n.url = url
	n.maxCapacity = maxCapacity
	n.running = make([]string, 0, maxCapacity)
	n.distinctFuncsLog = make([]int, 0, 1000)
	n.warmColdLogs = make([]bool, 0, 1000)

	return n
}

func (node *Node) setFinished(functionName string) {
	for i, name := range node.running {
		if name == functionName {
			node.running = append(node.running[:i], node.running[i+1:]...)
			return
		}
	}
}

var nodes [](*Node)

func main() {
	file, err := os.Open("data/events.csv")
	defer file.Close()

	if err != nil {
		panic(err)
	}

	nodes = make([]*Node, 0, 2)
	for i := 0; i < 2; i++ {
		nodes = append(nodes, newNode(i, "http://localhost:5000", 8))
	}

	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()

	i := 0
	for tick := 0; tick <= 60000; {
		logDistinctFunctionsCounts()

		if i >= len(rows) {
			break
		}
		row := rows[i]
		startTime, _ := strconv.Atoi(row[1])

		if tick >= startTime {
			name := row[0]
			node := leastLoaded()

			if name == "W1" || name == "W2" || name == "W3" {
				go runFunction(node, name)
			}
			i += 1

		} else {
			tick += 10
			time.Sleep(time.Second / 3)
		}
	}

	for i, node := range nodes {
		sum := 0
		for _, n := range node.distinctFuncsLog {
			sum += n
		}
		fmt.Printf("node %d: %.2f\n", i, float32(sum)/float32(len(node.distinctFuncsLog)))
	}
}

func logDistinctFunctionsCounts() {
	for _, n := range nodes {
		distincts := make(map[string]bool)
		for _, f := range n.running {
			distincts[f] = true
		}
		n.distinctFuncsLog = append(n.distinctFuncsLog, len(distincts))
	}
}
