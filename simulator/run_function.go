package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Result                FunctionExecutionResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
	LoadBalancingInfo     LoadBalancingInfoType
}
type FunctionExecutionResult struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}
type LoadBalancingInfoType struct {
	WorkerNodeId   int
	WorkerNodeUrl  string
	Algorithm      string
}
func runFunction(hostUrl string, functionName string) (Response, error) {
	resp, err := http.Get(hostUrl + "/execute?name=" + functionName)
	if err != nil {
		return Response{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	resp.Body.Close()

	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		println(err.Error())
		println(string(data))
	}
	return r, nil
}
