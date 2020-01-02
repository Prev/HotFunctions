package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Node struct {
	id          int
	url         string
	maxCapacity int
	// running     int64
	running map[string]int
	mutex   *sync.Mutex
}

func newNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.id = id
	n.url = url
	n.maxCapacity = maxCapacity
	// n.running = 0
	n.running = make(map[string]int)
	n.mutex = new(sync.Mutex)

	return n
}

func (node *Node) capacity() int {
	running := 0
	for _, val := range node.running {
		running += val
	}
	return node.maxCapacity - running
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

func (node *Node) runFunction(functionName string, onComplete func(string, int64)) {
	fmt.Printf("[run] %s at %d\n", functionName, node.id)

	startTime := time.Now().UnixNano()
	// atomic.AddInt64(&node.running, 1)

	node.mutex.Lock()
	if _, exists := node.running[functionName]; exists == false {
		node.running[functionName] = 0
	}
	node.running[functionName] += 1
	node.mutex.Unlock()

	resp, err := http.Get(node.url + "/execute?name=" + functionName)
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
	// atomic.AddInt64(&node.running, -1)
	node.mutex.Lock()
	node.running[functionName]--
	if node.running[functionName] == 0 {
		delete(node.running, functionName)
	}
	node.mutex.Unlock()

	onComplete(functionName, duration)
}
