package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	dtypes "github.com/Prev/HotFunctions/worker_front/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

const CPUQuota = 0

type Container struct {
	Name           string
	FunctionName   string
	IsRestMode     bool
	RestModePort   string
}

func CreateContainer(image Image) (Container, error) {
	ctx := context.Background()

	if image.IsRestMode {
		cont := Container{
			fmt.Sprintf("hf_%s__rest", strings.ToLower(image.FunctionName)),
			image.FunctionName,
			true,
			strconv.Itoa(RestModePortNo(image.FunctionName)),
		}
		logger.Printf("Start rest-mode container %s with port :%s\n", cont.Name, cont.RestModePort)

		_, err := cli.ContainerCreate(ctx,
			&container.Config{
				Image: image.Name,
				ExposedPorts: nat.PortSet{
					"8080/tcp": struct{}{},
				},
			},
			&container.HostConfig{
				Resources: container.Resources{ CPUQuota: CPUQuota },
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
			cont.Name,
		)
		if err != nil {
			if strings.Contains(err.Error(), "Conflict") {
				// Container exists
				return cont, nil
			}
			return Container{}, err
		}
		if err := cli.ContainerStart(ctx, cont.Name, types.ContainerStartOptions{}); err != nil {
			return Container{}, err
		}
		// Booting the container takes some time
		time.Sleep(time.Second / 2)

		return cont, nil

	} else {
		containerName := fmt.Sprintf("hf_%s__%d_%d",
			strings.ToLower(image.FunctionName),
			time.Now().Unix() % 100000,
			rand.Intn(100000),
		)

		cont := Container{
			containerName,
			image.FunctionName,
			false,
			"",
		}

		_, err := cli.ContainerCreate(
			ctx,
			&container.Config{
				Image: image.Name,
			},
			&container.HostConfig{
				Resources: container.Resources{ CPUQuota: CPUQuota },
			},
			nil,
			containerName,
		)
		if err != nil {
			return Container{}, err
		}

		return cont, nil
	}
}

func (c *Container) Run() (*dtypes.ContainerResponse, error) {
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

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		data = string(bytes)

		if len(data) < 10 {
			return nil, errors.New("rest container returns invalid result")
		}
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

func (c *Container) Remove() {
	ctx := context.Background()
	if err := cli.ContainerRemove(ctx, c.Name, types.ContainerRemoveOptions{Force: true}); err != nil {
		logger.Println(err.Error())
	}
}

func RestModePortNo(functionName string) int {
	b := md5.Sum([]byte(functionName))
	d := int(b[0]) + int(b[1]) + int(b[2])
	return 9000 + (d % 1000)
}