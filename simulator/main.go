package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: go run *.go <lb_address> [<event_stream_file>]")
		os.Exit(-1)
	}

	lbUrl := os.Args[1]
	eventStreamFileName := "data/events.csv"
	if len(os.Args) >= 3 {
		eventStreamFileName = os.Args[2]
	}

	logger := initLogger()

	// Events are predetermined and described at data/events.csv
	// Simulator run functions at specific time, using `time.Sleep` method.
	// The callback function will be executed based on startTime of the `events.csv` file,
	// and executed with goroutine, which means multiple functions can be run concurrently
	startSimulation(eventStreamFileName, func(functionName string, virtualTime int) {
		// StartTime will be recorded before running function & running scheduling algorithm
		startTime := time.Now().UnixNano()

		// Request function to load balancer
		resp, err := runFunction(lbUrl, functionName)
		if err != nil {
			panic(err)
		}

		// Calculate times from endTime and startTime
		endTime := time.Now().UnixNano()
		duration := (endTime - startTime) / int64(time.Millisecond)

		fmt.Printf("%s in %dms\n", functionName, duration)

		// Log result to the file
		logMsg := fmt.Sprintf(resp)
		logger.Output(2, logMsg)
	})

	time.Sleep(time.Second * 20)
}

func runFunction(hostUrl string, functionName string) (string, error) {
	resp, err := http.Get(hostUrl + "/execute?name=" + functionName)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	return string(data), nil
}
