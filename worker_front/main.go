package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

type lambdaResponse struct {
	StatusCode int64  `json:"statusCode"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	Body       string `json:"body"`
}

func (h *frontHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]

	w.Header().Add("Content-Type", "application/json")

	if len(nameParam) == 0 {
		resp := failResponse{false, "param 'name' is not given"}
		bytes, _ := json.Marshal(resp)
		w.Write(bytes)
		return
	}

	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{})

	for _, image := range images {
		splitted := strings.Split(image.RepoTags[0], ":")
		imageName := splitted[0]

		if imageName == nameParam[0] {
			startTime := makeTimestamp()
			log := runContainer(imageName)
			endTime := makeTimestamp()

			// Magic strings
			tmp := strings.Split(log, "-----")
			jsonLength, _ := strconv.ParseInt(tmp[1], 10, 64)
			idx := strings.Index(log, "-----=")
			jsonStr := log[idx+6 : idx+6+int(jsonLength)]

			var lr lambdaResponse
			err := json.Unmarshal([]byte(jsonStr), &lr)

			if err != nil {
				panic(err)
			}

			resp := successResponse{
				lr.Body,
				true,
				endTime - startTime,
				lr.EndTime - lr.StartTime,
			}
			bytes, _ := json.Marshal(resp)
			w.Write(bytes)
			return
		}
	}

	// TODO: automatic build image

	resp := failResponse{false, "fail to find image " + nameParam[0]}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)
}

func runContainer(imageName string) string {
	ctx := context.Background()

	containerName := fmt.Sprintf("%s_%d", imageName, rand.Intn(10000))

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, containerName)

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	return buf.String()
}
