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

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
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
	Result                string
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func (h *frontHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]

	w.Header().Add("Content-Type", "application/json")

	if len(nameParam) == 0 {
		resp := failResponse{true, "param 'name' is not given"}
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		return
	}

	startTime := makeTimestamp()

	functionName := nameParam[0]
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
		out.Body,
		imageFound,
		endTime - startTime,
		out.EndTime - out.StartTime,
	}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)
}
