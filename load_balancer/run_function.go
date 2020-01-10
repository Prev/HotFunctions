package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Prev/LALB/load_balancer/scheduler"
)

type WorkerResponse struct {
	Result                FunctionExecutionResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}
type FunctionExecutionResult struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func runFunction(node *scheduler.Node, functionName string, onComplete func(string, int64)) {
	fmt.Printf("[run] %s at %d\n", functionName, node.Id)

	startTime := time.Now().UnixNano()
	resp, err := http.Get(node.Url + "/execute?name=" + functionName)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	var result WorkerResponse
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

	log := fmt.Sprintf("%d %s %s %d %d %d", node.Id, functionName, startType, startTime/int64(time.Millisecond), duration, latency)
	logger.Output(2, log)

	onComplete(functionName, duration)
}
