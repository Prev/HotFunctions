package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dtypes "github.com/Prev/HotFunctions/worker_front/types"
)

type Container struct {
	Name           string
	FunctionName   string
	Reusable       bool
	IsRestMode     bool
	RestModePort   string
}

func CreateContainer(image Image) (Container, error) {
	ctx := context.Background()

	containerName := fmt.Sprintf("%s_%s__%d_%d",
		"hf_",
		strings.ToLower(image.FunctionName),
		time.Now().Unix() % 100000,
		rand.Intn(100000),
	)

	cont := Container{
		containerName,
		image.FunctionName,
		false,
		false,
		"",
	}
	if image.IsRestMode {
		cont.Reusable = true
		cont.IsRestMode = true
		cont.RestModePort = strconv.Itoa(rand.Intn(1000) + 9000)

		logger.Println("Start rest-mode container with port ", cont.RestModePort)

		_, err := cli.ContainerCreate(ctx,
			&container.Config{
				Image: image.Name,
				ExposedPorts: nat.PortSet{
					"8080/tcp": struct{}{},
				},
			},
			&container.HostConfig{
				PortBindings: nat.PortMap{
					"8080/tcp": []nat.PortBinding{
						{
							HostIP:   "127.0.0.1",
							HostPort: cont.RestModePort,
						},
					},
				},
			},
			nil,
			containerName,
		)
		if err != nil {
			return Container{}, err
		}
		if err := cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
			return Container{}, err
		}
		time.Sleep(time.Second)

	} else {
		_, err := cli.ContainerCreate(ctx, &container.Config{
			Image: image.Name,
		}, nil, nil, containerName)
		if err != nil {
			return Container{}, err
		}
	}

	return cont, nil
}

func containerBelongsToFunction(containerName string, functionName string) bool {
	return strings.Split(containerName, "__")[0] == "hf_" + strings.ToLower(functionName)
}

func (c Container) Run() (*dtypes.ContainerResponse, error) {
	containerID := c.Name
	ctx := context.Background()

	data := ""

	if c.IsRestMode == false {
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
		data = buf.String()

	} else {
		resp, err := http.Get("http://localhost:" + c.RestModePort)
		if err != nil {
			return nil, err
		}

		bytes, _ := ioutil.ReadAll(resp.Body)
		data = string(bytes)
	}

	// Handle magic strings
	// Log can be multiple lines if the container is executed multiple times.
	// So we use the last line of the logs, and the line separator is `==--==--==--==--==` in our system.
	lines := strings.Split(data, "==--==--==--==--==")
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

	var fr dtypes.ContainerResponse
	err := json.Unmarshal([]byte(jsonStr), &fr)

	if err != nil {
		return nil, err
	}

	return &fr, nil
}
