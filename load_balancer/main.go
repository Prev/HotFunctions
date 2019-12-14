package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	id          int
	ip          string
	maxCapacity int
	running     []string
	endTimes    []int
}

func newNode(id int, ip string, maxCapacity int) *Node {
	n := new(Node)
	n.id = id
	n.ip = ip
	n.maxCapacity = maxCapacity
	n.running = make([]string, 0, maxCapacity)
	n.endTimes = make([]int, 0, maxCapacity)

	return n
}

var nodes [](*Node)

func getFinishedFunctions(node *Node, currentTimestamp int) []string {
	finished := make([]string, 0, len(node.running))

	for i, time := range node.endTimes {
		if currentTimestamp >= time {
			finished = append(finished, node.running[i])
		}
	}

	return finished
}

func setFinished(node *Node, functionName string) {
	// println(functionName, "is finished")
	for i, name := range node.running {
		if name == functionName {
			node.running = append(node.running[:i], node.running[i+1:]...)
			node.endTimes = append(node.endTimes[:i], node.endTimes[i+1:]...)
			return
		}
	}
}

func main() {
	file, err := os.Open("data/events.csv")
	defer file.Close()

	if err != nil {
		panic(err)
	}

	nodes = make([]*Node, 0, 8)
	var distinctsPerNode [8][]int

	for i := 0; i < 8; i++ {
		nodes = append(nodes, newNode(i, "", 8))
		distinctsPerNode[i] = make([]int, 0, 60000)
	}

	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()

	i := 0
	for time := 0; time <= 60000; {
		for ni, n := range nodes {
			for _, fn := range getFinishedFunctions(n, time) {
				setFinished(n, fn)
			}
			distincts := make(map[string]bool)
			for _, f := range n.running {
				distincts[f] = true
			}
			distinctsPerNode[ni] = append(distinctsPerNode[ni], len(distincts))
		}

		if i >= len(rows) {
			break
		}
		row := rows[i]
		startTime, _ := strconv.Atoi(row[1])

		if time == startTime {
			node := leastLoaded()

			name := row[0]
			duration, _ := strconv.Atoi(row[2])
			endTime := startTime + duration

			node.running = append(node.running, name)
			node.endTimes = append(node.endTimes, endTime)

			// println(name, "is running at", node.id)
			i += 1

		} else {
			time += 1
		}
	}

	for i, dn := range distinctsPerNode {
		sum := 0
		for _, n := range dn {
			sum += n
		}
		fmt.Printf("node %d: %.2f\n", i, float32(sum)/float32(len(dn)))
	}
}
