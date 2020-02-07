package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger *log.Logger

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

	logger = initLogger()

	// Events are predetermined and described at data/events.csv
	// Simulator run functions at specific time, using `time.Sleep` method.
	// The callback function will be executed based on startTime of the `events.csv` file,
	// and executed with goroutine, which means multiple functions can be run concurrently
	startSimulation(eventStreamFileName, func(functionName string, virtualTime int) {
		// StartTime will be recorded before running function & running scheduling algorithm
		startTime := time.Now().UnixNano()

		// Request function call to load balancer
		resp, err := runFunction(lbUrl, functionName)
		if err != nil {
			panic(err)
		}

		// Calculate times from endTime and startTime
		endTime := time.Now().UnixNano()
		duration := (endTime - startTime) / int64(time.Millisecond)

		// InternalExecutionTime is pure execution time of the user function
		// Latency consists of Docker build time, container running time, scheduling algorithm execution time, network latency, etc.
		latency := duration - resp.InternalExecutionTime

		startType := "cold"
		if resp.IsWarm {
			startType = "warm"
		}

		nodeId := resp.LoadBalancingInfo.WorkerNodeId

		// Print to stdout
		fmt.Printf("%s in %dms, at node %d, %s start, latency %dms - %s\n", functionName, duration, nodeId, startType, latency, resp.Result.Body)

		// Log result to the file
		logMsg := fmt.Sprintf("%d %s %s %d %d %d", nodeId, functionName, startType, startTime/int64(time.Millisecond), duration, latency)
		logger.Output(2, logMsg)
	})

	time.Sleep(time.Second * 20)
}

