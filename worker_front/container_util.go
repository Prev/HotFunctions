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

func CreateContainer(functionName string) (string, error) {
	imageName := imageTagName(functionName)
	ctx := context.Background()

	containerName := fmt.Sprintf("%s__%d_%d", imageName, time.Now().Unix(), rand.Intn(100000))

	_, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, containerName)

	if err != nil {
		return "", err
	}

	//return resp.ID, nil
	return containerName, nil
}

func containerBelongsToFunction(containerName string, functionName string) bool {
	return strings.Split(containerName, "__")[0] == imageTagName(functionName)
}

func RunContainer(containerID string) (*FunctionResponse, error) {
	ctx := context.Background()

	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	log := buf.String()

	// Handle magic strings
	// Log can be multiple lines if the container is executed multiple times.
	// So we use the last line of the logs, and the line separator is `==--==--==--==--==` in our system.
	lines := strings.Split(log, "==--==--==--==--==")
	lastLine := lines[len(lines) - 2]

	// Sample line of the log is like below:
	// -=-=-=-=-=125-=-=-=-=-=>{"startTime": 1581315321872, "endTime": 1581315322914, "result": {"statusCode": 200, "body": "{\"ret\": 2499621.732014556}"}}
	// To handle the log, get the length of the json string firstly.
	tmp := strings.Split(lastLine, "-=-=-=-=-=")[1]
	jsonLength, _ := strconv.ParseInt(tmp, 10, 64)

	// Then get the first index of the body.
	idx := strings.Index(lastLine, "-=-=-=-=-=>")
	// Finally slice the log with the length
	jsonStr := lastLine[idx+11 : idx+11+int(jsonLength)]

	var fr FunctionResponse
	err = json.Unmarshal([]byte(jsonStr), &fr)

	if err != nil {
		return nil, err
	}

	return &fr, nil
}
