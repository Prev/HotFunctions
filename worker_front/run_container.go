package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type FunctionResponse struct {
	StartTime int64                  `json:"startTime"`
	EndTime   int64                  `json:"endTime"`
	Result    FunctionResponseResult `json:"result"`
}

type FunctionResponseResult struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func RunContainer(imageName string) (*FunctionResponse, error) {
	ctx := context.Background()

	containerName := fmt.Sprintf("%s_%d_%d", imageName, time.Now().Unix(), rand.Intn(100000))

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, containerName)

	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{}); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	log := buf.String()

	// Magic strings
	tmp := strings.Split(log, "-----")
	jsonLength, _ := strconv.ParseInt(tmp[1], 10, 64)
	idx := strings.Index(log, "-----=")
	jsonStr := log[idx+6 : idx+6+int(jsonLength)]

	var fr FunctionResponse
	err = json.Unmarshal([]byte(jsonStr), &fr)

	if err != nil {
		return nil, err
	}

	return &fr, nil
}
