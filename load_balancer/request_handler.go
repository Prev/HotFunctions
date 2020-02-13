package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Prev/HotFunctions/worker_front/types"
)

type RequestHandler struct {
	http.Handler
}

func newRequestHandler() *RequestHandler {
	h := new(RequestHandler)
	return h
}

type FailResponse struct {
	Error   bool
	Message string
}

type Response struct {
	types.ExecSuccessResponse
	LoadBalancingInfo LoadBalancingInfoType
}

type LoadBalancingInfoType struct {
	WorkerNodeId   int
	WorkerNodeUrl  string
	Algorithm      string
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Println(req.Method, req.URL.String())
	w.Header().Add("Content-Type", "application/json")

	switch req.URL.Path {
	case "/execute":
		h.ExecFunction(&w, req)
	default:
		w.WriteHeader(404)
		writeFailResponse(&w, "404 Not found on given path")
	}
}

func (h *RequestHandler) ExecFunction(w *http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]

	if len(nameParam) == 0 {
		writeFailResponse(w, "param 'name' is not given")
		return
	}
	functionName := nameParam[0]

	node, err := sched.Select(functionName)
	if err != nil {
		writeFailResponse(w, "error on selecting node from scheduler")
		return
	}
	defer sched.Finished(node, functionName)

	ret := Response{}
	if fakeMode {
		ret.Result = types.ContainerResponseData{0, ""}
		ret.InternalExecutionTime = 500
		ret.ExecutionTime = 500
		ret.Meta = types.FunctionExecutionMetaData{false, false, false, "", ""}
		time.Sleep(time.Second)

	} else {
		resp, err := http.Get(node.Url + "/execute?name=" + functionName)
		if err != nil {
			writeFailResponse(w, "error on sending http request on worker node")
			return
		}

		bytes, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(bytes, &ret); err != nil {
			println(err.Error())
			println(string(bytes))
			writeFailResponse(w, "error on parsing json from worker node")
			return
		}
	}

	ret.LoadBalancingInfo = LoadBalancingInfoType{
		node.Id,
		node.Url,
		schedType,
	}
	bytes, _ := json.Marshal(ret)
	(*w).Write(bytes)
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}
