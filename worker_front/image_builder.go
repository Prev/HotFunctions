package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mholt/archiver"
)

const sampleFunctionsBucket = "lalb-sample-functions"
const downloadPathPrefix = "_downloads/"
const envPathPrefix = "envs/"

func buildImage(functionName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Download files of the function
	println("Download file from S3")
	functionPath, err := downloadFile(sess, functionName)
	if err != nil {
		return err
	}

	// Read config.json file
	println("Reading config...")
	config, err := getConfigOfTheFunction(functionPath)
	if err != nil {
		return err
	}

	println("Making tar file...")
	tarPath, err := makeTarFile(functionPath, config["environment"])
	if err != nil {
		return err
	}

	// Build docker image
	println("Build image...")
	if err := buildImageWithTar(functionName, tarPath); err != nil {
		return err
	}

	return nil
}

func downloadFile(sess *session.Session, functionName string) (string, error) {
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

func getConfigOfTheFunction(functionPath string) (map[string]string, error) {
	jsonFile, err := os.Open(functionPath + "/config.json")
	defer jsonFile.Close()

	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func makeTarFile(functionPath string, envType string) (string, error) {
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
