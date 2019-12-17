package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Node struct {
	id          int
	url         string
	maxCapacity int
	running     map[string]string
}

func newNode(id int, url string, maxCapacity int) *Node {
	n := new(Node)
	n.id = id
	n.url = url
	n.maxCapacity = maxCapacity
	n.running = make(map[string]string)

	return n
}

type workerResponse struct {
	Result                string
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func (node *Node) runFunction(functionName string) {
	fmt.Printf("[run] %s at %d\n", functionName, node.id)

	startTime := time.Now().UnixNano()
	uid := fmt.Sprintf("%s-%d-%d", functionName, startTime, rand.Intn(100000))
	node.running[uid] = functionName

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

	delete(node.running, uid)
}
