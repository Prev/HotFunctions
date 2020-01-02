package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/docker/api/types"
	"github.com/mholt/archiver"
)

const DOWNLOAD_PATH_PREFIX = "_downloads/"
const ENV_PATH_PREFIX = "envs/"

type ImageBuilder struct {
	isBuilding map[string]bool
	mutex      *sync.Mutex
}

func newImageBuilder() *ImageBuilder {
	b := new(ImageBuilder)
	b.isBuilding = make(map[string]bool)
	b.mutex = new(sync.Mutex)
	return b
}

func (b *ImageBuilder) BuildSafe(functionName string) error {
	b.mutex.Lock()
	if b.isBuilding[functionName] == true {
		// Wait until image is built
		logger.Printf("Image for function '%s' not found. Wait for build compeletion...\n", functionName)
		for {
			b.mutex.Unlock()
			time.Sleep(time.Second / 20)
			b.mutex.Lock()

			if b.isBuilding[functionName] == false {
				break
			}
		}
		b.mutex.Unlock()

	} else {
		b.isBuilding[functionName] = true
		b.mutex.Unlock()

		logger.Printf("Image for function '%s' not found. Image build start.\n", functionName)
		if err := b.Build(functionName); err != nil {
			return err
		}

		b.mutex.Lock()
		b.isBuilding[functionName] = false
		b.mutex.Unlock()

		logger.Printf("Image for function '%s' build fin.\n", functionName)
	}

	return nil
}

func (b *ImageBuilder) Build(functionName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Download files of the function
	functionPath, err := b.downloadFiles(sess, functionName)
	if err != nil {
		return err
	}

	// Read config.json file
	config, err := b.getConfigOfTheFunction(functionPath)
	if err != nil {
		return err
	}

	// Make tar file for docker
	tarPath, err := b.makeTarFile(functionPath, config["environment"])
	if err != nil {
		return err
	}

	// Build docker image
	if err := b.buildImageWithTar(functionName, tarPath); err != nil {
		return err
	}

	return nil
}

func (b *ImageBuilder) downloadFiles(sess *session.Session, functionName string) (string, error) {
	os.MkdirAll(DOWNLOAD_PATH_PREFIX, 0700)
	zipFilePath := DOWNLOAD_PATH_PREFIX + functionName + ".zip"
	destPath := DOWNLOAD_PATH_PREFIX + functionName

	// Remove old files
	os.RemoveAll(zipFilePath)
	os.RemoveAll(destPath)

	// Download zip file
	resp, err := http.Get(USER_FUNCTION_URL_PREFIX + functionName + ".zip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Unzip file
	z := archiver.Zip{}
	if err := z.Unarchive(zipFilePath, DOWNLOAD_PATH_PREFIX); err != nil {
		return "", err
	}
	os.RemoveAll(zipFilePath)

	return destPath, nil
}

func (b *ImageBuilder) getConfigOfTheFunction(functionPath string) (map[string]string, error) {
	jsonFile, err := os.Open(functionPath + "/config.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func (b *ImageBuilder) makeTarFile(functionPath string, envType string) (string, error) {
	tarFilePath := functionPath + ".tar"
	fileList := []string{}

	// Add env files
	envDir := ENV_PATH_PREFIX + envType
	entries, err := ioutil.ReadDir(envDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		fileList = append(fileList, envDir+"/"+entry.Name())
	}

	// Add downloaded files
	entries, err = ioutil.ReadDir(functionPath)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		fileList = append(fileList, functionPath+"/"+entry.Name())
	}

	// Remove old
	os.RemoveAll(tarFilePath)

	// Make a tar file
	t := archiver.Tar{}
	if err := t.Archive(fileList, tarFilePath); err != nil {
		return "", err
	}

	return tarFilePath, nil
}

func (b *ImageBuilder) buildImageWithTar(functionName string, tarPath string) error {
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
	if err != nil {
		return err
	}

	// Wait until bulid finished
	io.Copy(ioutil.Discard, out.Body)
	out.Body.Close()

	return nil
}
