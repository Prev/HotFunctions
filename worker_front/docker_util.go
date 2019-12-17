package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type lambdaResponse struct {
	StartTime int64                `json:"startTime"`
	EndTime   int64                `json:"endTime"`
	Result    lambdaResponseResult `json:"result"`
}

type lambdaResponseResult struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func imageTagName(functionName string) string {
	return "lalb_" + strings.ToLower(functionName)
}

func runContainer(imageName string) (*lambdaResponse, error) {
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
		Dockerfile: "/Dockerfile",
		Tags:       []string{imageTagName(functionName)},
	}

	out, err := cli.ImageBuild(ctx, dockerBuildContext, opt)
	//defer out.Body.Close()

	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, out.Body)
	out.Body.Close()

	//if _, err := io.Copy(os.Stdout, out.Body); err != nil {
	//	return err
	//}

	return nil
}

func removeOldImages() error {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	tmp := make([]string, 0, len(cachedImages))
	for key, val := range cachedImages {
		tmp = append(tmp, strconv.FormatInt(val, 10)+":"+key)
	}
	sort.Strings(tmp)

	for i := 0; i < len(tmp)-IMAGE_CACHE_NUM; i++ {
		spliited := strings.Split(tmp[i], ":")
		key := spliited[1]
		delete(cachedImages, key)
	}

	for _, image := range images {
		imageName := strings.Split(image.RepoTags[0], ":")[0]
		_, exists := cachedImages[imageName]

		if len(imageName) > 5 && imageName[0:5] == "lalb_" && exists == false {
			_, err := cli.ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{
				Force: true,
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}
