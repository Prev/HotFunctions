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
	Message string
	State   CachingOptions
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
	message := "nothing changed"

	imageLimit := q["image_limit"]
	if len(imageLimit) > 0 {
		val, _ := strconv.Atoi(imageLimit[0])
		h.functionRunner.cachingOptions.ImageLimit = val
		message = "configure changed"
	}

	containerPoolLimit := q["container_pool_limit"]
	if len(containerPoolLimit) > 0 {
		val, _ := strconv.Atoi(containerPoolLimit[0])
		h.functionRunner.cachingOptions.ContainerPoolLimit = val
		message = "configure changed"
	}

	containerPoolNum := q["container_pool_num"]
	if len(containerPoolNum) > 0 {
		val, _ := strconv.Atoi(containerPoolNum[0])
		h.functionRunner.cachingOptions.ContainerPoolNum = val
		message = "configure changed"
	}

	usingRestMode := q["using_rest_mode"]
	if len(usingRestMode) > 0 {
		if usingRestMode[0] == "true" {
			h.functionRunner.cachingOptions.UsingRestMode = true
			message = "configure changed"
		} else if usingRestMode[0] == "false" {
			h.functionRunner.cachingOptions.UsingRestMode = false
			message = "configure changed"
		}

		h.functionRunner.images = make(map[string]Image)
	}

	resp := ConfigureSuccessResponse{message, h.functionRunner.cachingOptions}
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
	if out == nil {
		writeFailResponse(w, "error on running a function")
		return
	}

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
