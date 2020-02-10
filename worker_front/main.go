package main // import "github.com/Prev/LALB/worker_front"

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/sevlyar/go-daemon"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

var cli *client.Client
var logger *log.Logger

var UserFunctionUrlPrefix string

type CachingOptions struct {
	ImageLimit         int
	ContainerPoolLimit int
	ContainerPoolNum   int
}

func main() {
	var err error

	port := 8222
	if len(os.Args) >= 2 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	if len(os.Args) >= 3 && os.Args[2] == "-d" {
		// Run as a daemon if second argument is "-d"
		cntxt := &daemon.Context{
			PidFileName: "daemon.pid",
			PidFilePerm: 0644,
			LogFileName: "daemon.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		child, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if child != nil {
			// Parent
			println("Start worker-front as a daemon")
		} else {
			// Child daemon
			defer cntxt.Release()
			runWorkerFront(port)
		}

	} else {
		runWorkerFront(port)
	}
}

func runWorkerFront(port int) {
	var err error

	UserFunctionUrlPrefix = getEnvString("USER_FUNCTION_URL_PREFIX", "https://lalb-sample-functions.s3.ap-northeast-2.amazonaws.com/")
	goMaxProcs := getEnvInt("GOMAXPROCS", 8)

	imageCacheLimit := getEnvInt("IMAGE_CACHE_LIMIT", -1)
	containerPoolLimit := getEnvInt("CONTAINER_POOL_LIMIT", -1)
	containerPoolNum := getEnvInt("CONTAINER_POOL_NUM", 4)

	runtime.GOMAXPROCS(goMaxProcs)

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	if cli, err = client.NewClientWithOpts(client.WithVersion("1.39")); err != nil {
		panic(err)
	}

	logger.Printf("GOMAXPROCS: %d\n", goMaxProcs)
	logger.Printf("Server listening at :%d\n", port)

	http.Handle("/", newRequestHandler(CachingOptions{
		imageCacheLimit,
		containerPoolLimit,
		containerPoolNum,
	}))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}
