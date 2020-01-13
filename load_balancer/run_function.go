package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

func runFunction(node *scheduler.Node, functionName string) (WorkerResponse, error) {
	// Since worker node provide RESTful API to execute a function,
	// Send Get request to it
	resp, err := http.Get(node.Url + "/execute?name=" + functionName)
	if err != nil {
		return WorkerResponse{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WorkerResponse{}, err
	}
	resp.Body.Close()

	var result WorkerResponse
	if err := json.Unmarshal(data, &result); err != nil {
		println(err.Error())
		println(string(data))
	}

	// Since this function will called as goroutine, send result by callback
	return result, nil
}
