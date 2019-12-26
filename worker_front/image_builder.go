package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/docker/docker/api/types"
	"github.com/mholt/archiver"
)

const sampleFunctionsBucket = "lalb-sample-functions"
const downloadPathPrefix = "_downloads/"
const envPathPrefix = "envs/"

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

func (b *ImageBuilder) Build2(functionName string) error {
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
	functionPath, err := b.downloadFile(sess, functionName)
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

func (b *ImageBuilder) downloadFile(sess *session.Session, functionName string) (string, error) {
	os.MkdirAll(downloadPathPrefix, 0700)
	zipFilePath := downloadPathPrefix + functionName + ".zip"
	destPath := downloadPathPrefix + functionName

	os.RemoveAll(zipFilePath)
	os.RemoveAll(destPath)

	file, err := os.Create(zipFilePath)
	defer file.Close()

	if err != nil {
		return "", err
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(sampleFunctionsBucket),
			Key:    aws.String(functionName + ".zip"),
		})

	if err != nil {
		return "", err
	}

	z := archiver.Zip{}
	if err := z.Unarchive(zipFilePath, downloadPathPrefix); err != nil {
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
	fileList := []string{}

	envDir := envPathPrefix + envType
	entries, err := ioutil.ReadDir(envDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		fileList = append(fileList, envDir+"/"+entry.Name())
	}

	entries, _ = ioutil.ReadDir(functionPath)
	for _, entry := range entries {
		fileList = append(fileList, functionPath+"/"+entry.Name())
	}
	tarFilePath := functionPath + ".tar"
	os.RemoveAll(tarFilePath)

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
	io.Copy(ioutil.Discard, out.Body)
	out.Body.Close()

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
