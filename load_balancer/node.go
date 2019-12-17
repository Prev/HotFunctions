package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Node struct {
	id          int
	url         string
	maxCapacity int
	running     int
}

func newNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.id = id
	n.url = url
	n.maxCapacity = maxCapacity
	n.running = 0

	return n
}

type workerResponse struct {
	Result                string
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func (node *Node) runFunction(functionName string, logger *log.Logger) {
	fmt.Printf("[run] %s at %d\n", functionName, node.id)

	startTime := time.Now().UnixNano()
	node.running++

	resp, err := http.Get(node.url + "?key=CS530&name=" + functionName)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result workerResponse
	json.Unmarshal(data, &result)

	endTime := time.Now().UnixNano()
	duration := (endTime - startTime) / int64(time.Millisecond)

	startType := "cold"
	if result.IsWarm {
		startType = "warm"
	}

	latency := duration - result.InternalExecutionTime

	fmt.Printf("[finished] %s in %dms, %s start, latency %dms\n",
		functionName, duration, startType, latency)

	log := fmt.Sprintf("%d %d %s %s %d %d", endTime/1000, node.id, functionName, startType, duration, latency)
	logger.Output(2, log)

	node.running--
}
