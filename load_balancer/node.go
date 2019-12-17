package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Node struct {
	id          int
	url         string
	maxCapacity int
	running     int64
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
	Result                lambdaResponseResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}
type lambdaResponseResult struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func (node *Node) runFunction(functionName string, logger *log.Logger) {
	fmt.Printf("[run] %s at %d\n", functionName, node.id)

	startTime := time.Now().UnixNano()
	atomic.AddInt64(&node.running, 1)

	resp, err := http.Get(node.url + "?key=CS530&name=" + functionName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	var result workerResponse
	if err := json.Unmarshal(data, &result); err != nil {
		println(err.Error())
		println(string(data))
	}

	if result.ExecutionTime == 0 {
		println(string(data))
	}

	endTime := time.Now().UnixNano()
	duration := (endTime - startTime) / int64(time.Millisecond)

	startType := "cold"
	if result.IsWarm {
		startType = "warm"
	}

	latency := duration - result.InternalExecutionTime

	fmt.Printf("[finished] %s in %dms, %s start, latency %dms - %d %d %s\n",
		functionName, duration, startType, latency, result.ExecutionTime, result.InternalExecutionTime, result.Result.Body)

	log := fmt.Sprintf("%d %s %s %d %d %d", node.id, functionName, startType, startTime/int64(time.Millisecond), duration, latency)
	logger.Output(2, log)

	// node.running--
	atomic.AddInt64(&node.running, -1)
}
