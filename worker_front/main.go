package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client
var cachedImages map[string]int64
var mutex *sync.Mutex

const IMAGE_CACHE_NUM = 2
const REQUEST_API_KEY = "CS530"

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	cachedImages = make(map[string]int64)
	mutex = new(sync.Mutex)

	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		panic(err)
	}

	port := 8222
	if len(os.Args) >= 2 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("server listening at :%d\n", port)
	http.Handle("/", new(FrontHandler))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}

type FrontHandler struct {
	http.Handler
}

type FailResponse struct {
	Error   bool
	Message string
}

type SuccessResponse struct {
	Result                lambdaResponseResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func (h *FrontHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]
	apiKeyParam := q["key"]

	w.Header().Add("Content-Type", "application/json")

	if len(apiKeyParam) == 0 || apiKeyParam[0] != REQUEST_API_KEY {
		writeFailResponse(&w, "param 'key' is not given or not valid")
		return
	}

	if len(nameParam) == 0 {
		writeFailResponse(&w, "param 'name' is not given")
		return
	}

	println("function requested:", nameParam[0])

	startTime := makeTimestamp()

	functionName := nameParam[0]
	targetImageName := imageTagName(functionName)

	imageFound := false
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		writeFailResponse(&w, err.Error())
		return
	}

	for _, image := range images {
		if len(image.RepoTags) > 0 {
			splitted := strings.Split(image.RepoTags[0], ":")

			if len(splitted) > 0 && splitted[0] == targetImageName {
				imageFound = true
				break
			}
		}
	}

	// build image if image not exist
	if imageFound == false {
		fmt.Printf("Image '%s' not found. Start to build\n", targetImageName)
		if err := buildImage(functionName); err != nil {
			writeFailResponse(&w, err.Error())
			return
		}
	}

	out, err := runContainer(targetImageName)
	if err != nil {
		writeFailResponse(&w, err.Error())
		return
	}

	endTime := makeTimestamp()
	println("fin", functionName)

	resp := SuccessResponse{
		out.Result,
		imageFound,
		endTime - startTime,
		out.EndTime - out.StartTime,
	}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)

	go handleCachedImages(targetImageName)
}

func handleCachedImages(imageName string) {
	mutex.Lock()
	cachedImages[imageName] = time.Now().Unix()

	if err := removeOldImages(); err != nil {
		//println(err.Error())
	}

	mutex.Unlock()
}
