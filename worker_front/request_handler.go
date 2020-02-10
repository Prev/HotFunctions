package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// RequestHandler of the worker front
type RequestHandler struct {
	http.Handler
	functionRunner *FunctionRunner
}

func newRequestHandler(options CachingOptions) *RequestHandler {
	h := new(RequestHandler)
	h.functionRunner = newFunctionRunner(options)
	return h
}

type FailResponse struct {
	Error   bool
	Message string
}

type ConfigureSuccessResponse struct {
	Result  string
	Message string
}

type ExecSuccessResponse struct {
	Result                FunctionResponseResult
	ExecutionTime         int64
	InternalExecutionTime int64
	Meta                  runningMetaData
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Println(req.Method, req.URL.String())
	w.Header().Add("Content-Type", "application/json")

	switch req.URL.Path {
	case "/configure":
		h.ConfigureWorker(&w, req)
	case "/execute":
		h.ExecFunction(&w, req)
	default:
		w.WriteHeader(404)
		writeFailResponse(&w, "404 Not found on given path")
	}
}

func (h *RequestHandler) ConfigureWorker(w *http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	cacheNumParam := q["cache_num"]

	if len(cacheNumParam) == 0 {
		writeFailResponse(w, "param 'cache_num' is not given")
		return
	}

	cacheNum, _ := strconv.Atoi(cacheNumParam[0])
	h.functionRunner.cachingOptions.imageCacheMaxNumber = cacheNum

	resp := ConfigureSuccessResponse{"success", "configure changed successfully"}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func (h *RequestHandler) ExecFunction(w *http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]

	if len(nameParam) == 0 {
		writeFailResponse(w, "param 'name' is not given")
		return
	}
	startTime := makeTimestamp()

	functionName := nameParam[0]
	out, meta := h.functionRunner.runFunction(functionName)

	endTime := makeTimestamp()
	logger.Println("fin", functionName)

	resp := ExecSuccessResponse{
		out.Result,
		endTime - startTime,
		out.EndTime - out.StartTime,
		meta,
	}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
