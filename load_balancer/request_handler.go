package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
		writeFailResponse(w, "err on selecting node from scheduler")
		return
	}

	resp, err := http.Get(node.Url + "/execute?name=" + functionName)
	sched.Finished(node, functionName)

	if err != nil {
		writeFailResponse(w, "err on sending http request on worker node")
		return
	}

	(*w).Header().Set("X-Balancing-Algorithm", schedType)
	(*w).Header().Set("X-Node-ID", strconv.Itoa(node.Id))
	(*w).Header().Set("X-Node-URL", node.Url)

	bytes, _ := ioutil.ReadAll(resp.Body)
	(*w).Write(bytes)
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}
