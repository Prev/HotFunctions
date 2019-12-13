package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type lambdaResponse struct {
	StatusCode int64  `json:"statusCode"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	Body       string `json:"body"`
}

func imageTagName(functionName string) string {
	return "lalb_" + strings.ToLower(functionName)
}

func runContainer(imageName string) (*lambdaResponse, error) {
	ctx := context.Background()

	containerName := fmt.Sprintf("%s_%d", imageName, rand.Intn(10000))

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

	var lr lambdaResponse
	err = json.Unmarshal([]byte(jsonStr), &lr)

	if err != nil {
		return nil, err
	}

	return &lr, nil
}

func buildImageWithTar(functionName string, tarPath string) error {
	dockerBuildContext, err := os.Open(tarPath)
	defer dockerBuildContext.Close()

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	opt := types.ImageBuildOptions{
		// Dockerfile: functionName + "/Dockerfile",
		Dockerfile: "/Dockerfile",
		Tags:       []string{imageTagName(functionName)},
	}

	out, err := cli.ImageBuild(ctx, dockerBuildContext, opt)
	defer out.Body.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, out.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
		return err
	}

	return nil
}
