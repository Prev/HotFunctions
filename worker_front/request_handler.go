package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

type RequestHandler struct {
	http.Handler
	cachedImages      map[string]int64
	cachedImagesMutex *sync.Mutex
	imageBuilder      *ImageBuilder
}

func newRequestHandler() *RequestHandler {
	h := new(RequestHandler)
	h.cachedImages = make(map[string]int64)
	h.cachedImagesMutex = new(sync.Mutex)
	h.imageBuilder = newImageBuilder()
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
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Println(req.Method, req.URL.String())
	w.Header().Add("Content-Type", "application/json")

	switch req.URL.Path {
	case "/execute":
		h.ExecFunction(&w, req)
	case "/configure":
		h.ConfigureWorker(&w, req)
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
	IMAGE_CACHE_NUM = cacheNum

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
	imageName := imageTagName(functionName)

	// logger.Println("function requested:", functionName)

	cacheExists := false
	h.cachedImagesMutex.Lock()
	if _, exists := h.cachedImages[imageName]; exists == true {
		cacheExists = true
	}
	h.cachedImagesMutex.Unlock()

BUILD_IMAGE:
	if cacheExists == false {
		h.imageBuilder.BuildSafe(functionName)
	}

	out, err := RunContainer(imageName)
	if err != nil {
		logger.Println(err.Error(), "Retry...")
		cacheExists = false
		goto BUILD_IMAGE
	}

	endTime := makeTimestamp()
	logger.Println("fin", functionName)

	resp := ExecSuccessResponse{
		out.Result,
		cacheExists,
		endTime - startTime,
		out.EndTime - out.StartTime,
	}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)

	go h.updateCachedImages(imageName)
}

func (h *RequestHandler) updateCachedImages(imageName string) {
	h.cachedImagesMutex.Lock()
	h.cachedImages[imageName] = time.Now().Unix()

	// Remove old images
	tmp := make([]string, 0, len(h.cachedImages))
	for key, val := range h.cachedImages {
		tmp = append(tmp, strconv.FormatInt(val, 10)+":"+key)
	}
	sort.Strings(tmp)

	for i := 0; i < len(tmp)-IMAGE_CACHE_NUM; i++ {
		spliited := strings.Split(tmp[i], ":")
		key := spliited[1]

		_, err := cli.ImageRemove(context.Background(), key, types.ImageRemoveOptions{Force: true})
		if err != nil {
			logger.Println(err.Error())
		}
		delete(h.cachedImages, key)
	}

	h.cachedImagesMutex.Unlock()
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
