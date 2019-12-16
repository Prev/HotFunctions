package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client
var cachedImages map[string]int64

const IMAGE_CACHE_NUM = 2
const REQUEST_API_KEY = "CS530"

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	// imageLastUsedTime = make(map[string]int64)
	cachedImages = make(map[string]int64)

	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.40"))

	if err != nil {
		panic(err)
	}

	println("server listening at :5000")
	http.Handle("/", new(frontHandler))
	http.ListenAndServe(":5000", nil)
}

type frontHandler struct {
	http.Handler
}

type failResponse struct {
	Error   bool
	Message string
}

type successResponse struct {
	Result                lambdaResponseResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func (h *frontHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]
	apiKeyParam := q["key"]

	w.Header().Add("Content-Type", "application/json")

	if len(apiKeyParam) == 0 || apiKeyParam[0] != REQUEST_API_KEY {
		resp := failResponse{true, "param 'key' is not given or not valid"}
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		return
	}

	if len(nameParam) == 0 {
		resp := failResponse{true, "param 'name' is not given"}
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		return
	}

	startTime := makeTimestamp()

	functionName := nameParam[0]
	println("function requested:", functionName)

	targetImageName := imageTagName(functionName)

	imageFound := false
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{})

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
		if err := buildImage(functionName); err != nil {
			resp := failResponse{true, err.Error()}
			bytes, _ := json.Marshal(resp)
			w.Write(bytes)
			return
		}
	}

	out, err := runContainer(targetImageName)

	if err != nil {
		resp := failResponse{true, err.Error()}
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		return
	}

	endTime := makeTimestamp()

	resp := successResponse{
		out.Result,
		imageFound,
		endTime - startTime,
		out.EndTime - out.StartTime,
	}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)

	// imageLastUsedTime[targetImageName] = time.Now().Unix()
	cachedImages[targetImageName] = time.Now().Unix()

	if err := removeOldImages(); err != nil {
		println(err)
	}
}
