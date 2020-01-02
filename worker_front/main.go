package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
)

var cli *client.Client
var logger *log.Logger

var IMAGE_CACHE_NUM int
var USER_FUNCTION_URL_PREFIX string

func imageTagName(functionName string) string {
	return "lalb_" + strings.ToLower(functionName)
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	GOMAXPROCS := getEnvInt("GOMAXPROCS", 8)
	IMAGE_CACHE_NUM = getEnvInt("IMAGE_CACHE_NUM", 4)
	USER_FUNCTION_URL_PREFIX = getEnvString("USER_FUNCTION_URL_PREFIX", "https://lalb-sample-functions.s3.ap-northeast-2.amazonaws.com/")

	runtime.GOMAXPROCS(GOMAXPROCS)

	var err error
	port := 8222
	if len(os.Args) >= 2 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	cli, err = client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("GOMAXPROCS: %d, IMAGE_CACHE_NUM: %d\n", GOMAXPROCS, IMAGE_CACHE_NUM)
	logger.Printf("Server listening at :%d\n", port)

	http.Handle("/", newRequestHandler())
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}
