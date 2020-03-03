package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const daemonPidName = "daemon.pid"
var cli *client.Client
var logger *log.Logger

var UserFunctionUrlPrefix string

type CachingOptions struct {
	UserCodeSizeLimit     int64
	ContainerPoolLimit    int
	ContainerPoolNum      int
	UsingRestMode         bool
	RestContainerLifeTime int
}

func main() {
	var err error

	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	switch os.Args[1] {
	case "start":
		port := 8222
		if len(os.Args) >= 3 {
			port, err = strconv.Atoi(os.Args[2])
			if err != nil {
				panic(err)
			}
		}
		if len(os.Args) >= 4 && os.Args[3] == "-d" {
			// Run as a daemon if second argument is "-d"
			cntxt := &daemon.Context{
				PidFileName: daemonPidName,
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

	case "stop":
		pid, err := ioutil.ReadFile(daemonPidName)
		if err != nil {
			panic(err)
		}

		_, err = exec.Command("kill", string(pid)).Output()
		if err != nil {
			panic(err)
		}

		if err := os.Remove(daemonPidName); err != nil {
			panic(err)
		}

		println("Daemon stopped")

	default:
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	println("Usage: \n" +
		"\tgo run *.go start\n" +
		"\tgo run *.go start [port]\n" +
		"\tgo run *.go start [port] [-d]\n" +
		"\tgo run *.go stop")
	os.Exit(-1)
}

func runWorkerFront(port int) {
	var err error

	UserFunctionUrlPrefix = getEnvString("USER_FUNCTION_URL_PREFIX", "https://lalb-sample-functions.s3.ap-northeast-2.amazonaws.com/")
	goMaxProcs := getEnvInt("GOMAXPROCS", 8)

	cachingOptions := CachingOptions{
		UserCodeSizeLimit:     getEnvInt64("USER_CODE_SIZE_LIMIT", -1),
		ContainerPoolLimit:    getEnvInt("CONTAINER_POOL_LIMIT", -1),
		ContainerPoolNum:      getEnvInt("CONTAINER_POOL_NUM", 6),
		UsingRestMode:         false,
		RestContainerLifeTime: getEnvInt("REST_CONTAINER_LIFE_TIME", 10),
	}
	if getEnvString("USING_REST_MODE", "") != "" {
		cachingOptions.UsingRestMode = true
	}

	runtime.GOMAXPROCS(goMaxProcs)

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	if cli, err = client.NewClientWithOpts(client.WithVersion("1.39")); err != nil {
		panic(err)
	}

	logger.Printf("GOMAXPROCS: %d\n", goMaxProcs)
	logger.Printf("Server listening at :%d\n", port)

	// Clean containers before running
	CleanContainers()

	http.Handle("/", newRequestHandler(cachingOptions))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}
