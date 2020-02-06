package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type LoadBalancerResponse struct {
	NodeId int
	NodeUrl string
}
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

func runFunction(hostUrl string, functionName string) (LoadBalancerResponse, WorkerResponse, error) {
	resp, err := http.Get(hostUrl + "/execute?name=" + functionName)
	if err != nil {
		return LoadBalancerResponse{}, WorkerResponse{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LoadBalancerResponse{}, WorkerResponse{}, err
	}
	resp.Body.Close()

	nodeId, _ := strconv.Atoi(resp.Header.Get("X-Node-Id"))
	nodeUrl := resp.Header.Get("X-Node-Url")

	lbr := LoadBalancerResponse{
		nodeId,
		nodeUrl,
	}

	var wr WorkerResponse
	if err := json.Unmarshal(data, &wr); err != nil {
		println(err.Error())
		println(string(data))
	}

	return lbr, wr, nil
}
