package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type workerResponse struct {
	Result                string
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func runFunction(node *Node, functionName string) {
	fmt.Printf("[run] %s at %d\n", functionName, node.id)

	startTime := time.Now().UnixNano()
	node.running = append(node.running, functionName)

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

	fmt.Printf("[finished]: %s in %dms, %s start, latency %dms\n",
		functionName, duration, startType, latency)

	node.setFinished(functionName)
}
