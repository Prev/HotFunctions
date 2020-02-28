package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
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
	WorkerNodeId     int
	WorkerNodeUrl    string
	Algorithm        string
	AlgorithmLatency float64
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Println(req.Method, req.URL.String())
	w.Header().Add("Content-Type", "application/json")

	switch req.URL.Path {
	case "/configure":
		h.ConfigureBalancer(&w, req)
	case "/clear":
		h.Clear(&w, req)
	case "/prepare":
		h.Prepare(&w, req)
	case "/execute":
		h.ExecFunction(&w, req)
	default:
		w.WriteHeader(404)
		writeFailResponse(&w, "404 Not found on given path")
	}
}

func (h *RequestHandler) Clear(w *http.ResponseWriter, req *http.Request) {
	wg := sync.WaitGroup{}
	for _, node := range nodes {
		wg.Add(1)
		go func(nodeUrl string) {
			defer wg.Done()

			_, err := http.Get( + "/clear")
			if err != nil {
				println(err.Error())
				(*w).Write([]byte("\"error\""))
			}
		}(node.Url)
	}
	wg.Wait()
	(*w).Write([]byte("done"))
}

func (h *RequestHandler) Prepare(w *http.ResponseWriter, req *http.Request) {
	wg := sync.WaitGroup{}
	for _, node := range nodes {
		wg.Add(1)

		go func(nodeUrl string) {
			defer wg.Done()
			_, err := http.Get(nodeUrl + "/prepare")
			if err != nil {
				println(err.Error())
				(*w).Write([]byte("\"error\""))
			}
		}(node.Url)
	}
	wg.Wait()
	(*w).Write([]byte("done"))
}

func (h *RequestHandler) ConfigureBalancer(w *http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if v := q["lb"]; len(v) > 0 {
		if err := setScheduler(v[0]); err != nil {
			println(err.Error())
		}
	}

	(*w).Write([]byte("{"))
	(*w).Write([]byte("\"lb\":"))
	(*w).Write([]byte("\"" + schedType + "\""))
	(*w).Write([]byte("\n"))

	for _, node := range nodes {
		(*w).Write([]byte(",\"" + node.Url +  "\":"))

		resp, err := http.Get(node.Url + "/configure?" + req.URL.RawQuery)
		if err != nil {
			println(err.Error())
			(*w).Write([]byte("\"error\""))

		} else {
			bytes, _ := ioutil.ReadAll(resp.Body)
			(*w).Write(bytes)
		}

		(*w).Write([]byte("\n"))
	}
	(*w).Write([]byte("}"))
}

func (h *RequestHandler) ExecFunction(w *http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]

	if len(nameParam) == 0 {
		writeFailResponse(w, "param 'name' is not given")
		return
	}
	functionName := nameParam[0]

	algorithmStartTime := float64(time.Now().UnixNano()) / float64(time.Millisecond)
	node, err := sched.Select(functionName)
	algorithmEndTime := float64(time.Now().UnixNano()) / float64(time.Millisecond)

	if err != nil {
		writeFailResponse(w, "error on selecting node from scheduler")
		return
	}
	defer sched.Finished(node, functionName)

	ret := Response{}
	if fakeMode {
		ret.Result = types.ContainerResponseData{0, ""}
		ret.InternalExecutionTime = 1000
		ret.ExecutionTime = 4000
		ret.Meta = types.FunctionExecutionMetaData{false, false, false, "", ""}
		time.Sleep(4 * time.Second)

	} else {
		resp, err := http.Get(node.Url + "/execute?name=" + functionName)
		if err != nil {
			println(err.Error())
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
		algorithmEndTime - algorithmStartTime,
	}
	bytes, _ := json.Marshal(ret)
	(*w).Write(bytes)
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}
