package main

import (
	"encoding/json"
	"fmt"
	"github.com/Prev/LALB/load_balancer/scheduler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type NodeConfigData struct {
	Url         string `json:"url"`
	MaxCapacity int    `json:"maxCapacity"`
}

var logger *log.Logger
var sched scheduler.Scheduler
var schedType string

func main() {
	if len(os.Args) != 2 {
		println("Usage: go run *.go ll|hash|ours|pasch")
		os.Exit(-1)
	}
	schedType = os.Args[1]

	// Modify nodes.config.json to configure worker nodes
	nodes := initNodesFromConfig("nodes.config.json")
	//nodes := make([]*scheduler.Node, 6)
	//for i := 0; i < 6; i++ {
	//	nodes[i] = scheduler.NewNode(i, "")
	//}

	switch schedType {
	case "ll":
		// Least Loaded Scheduler
		// Scheduler picks the node who has minimum executing tasks
		println("Using Least Loaded Scheduler")
		sched = scheduler.NewLeastLoadedScheduler(&nodes)

	case "hash":
		// Consistent Hashing Scheduler
		// Scheduler picks the node by Consistent Hashing algorithm where key is the function name
		println("Using Consistent Hashing Scheduler")
		sched = scheduler.NewConsistentHashingScheduler(&nodes, 8)

	case "pasch":
		// Consistent Hashing Scheduler
		// Scheduler picks the node by Consistent Hashing algorithm where key is the function name
		println("Using PASch Extended Scheduler")
		sched = scheduler.NewPASchExtendedScheduler(&nodes, 8)

	case "ours":
		// Proposing Greedy Scheduler
		println("Using Our Scheduler")
		sched = scheduler.NewOurScheduler(&nodes, 8, 6, 3)
	}

	port := 8111

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Server listening at :%d\n", port)

	http.Handle("/", newRequestHandler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}

	/*
	// Events are predetermined and described at data/events.csv
	// Simulator run functions at specific time, using `time.Sleep` method.
	// The callback function will be executed based on startTime of the `events.csv` file,
	// and executed with goroutine, which means multiple functions can be run concurrently
	simulator := simulator2.newSimulator("data/events.csv")
	simulator.Start(func(functionName string, virtualTime int) {
		// StartTime will be recorded before running function & running scheduling algorithm
		startTime := time.Now().UnixNano()

		// Select proper worker node by scheduler
		node, err := sched.Select(functionName)

		if err != nil {
			fmt.Printf("[fail running] %s: %s\n", functionName, err.Error())
			return
		}

		fmt.Printf("[run] %s at %d\n", functionName, node.Id)

		// 1~3 sec
		rn := rand.Float64() * 2 + 1
		time.Sleep(time.Second * time.Duration(rn))

		endTime := time.Now().UnixNano()
		duration := (endTime - startTime) / int64(time.Millisecond)
		latency := 0
		startType := "cold"

		// Log result to the file
		logMsg := fmt.Sprintf("%d %s %s %d %d %d", node.Id, functionName, startType, startTime/int64(time.Millisecond), duration, latency)
		logger.Output(2, logMsg)


		// Notify scheduler that the function execution is completed
		// Scheduler will modify its data structure (e.g. capacity-table) if it needs
		sched.Finished(node, functionName, duration)
	})

	time.Sleep(time.Second * 20)
	*/
}

// Init node list from the config file
func initNodesFromConfig(configFilePath string) []*scheduler.Node {
	nodeConfigFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer nodeConfigFile.Close()

	var nodeConfigs []NodeConfigData
	byteValue, _ := ioutil.ReadAll(nodeConfigFile)
	json.Unmarshal([]byte(byteValue), &nodeConfigs)

	nodes := make([]*scheduler.Node, len(nodeConfigs))
	for i, nc := range nodeConfigs {
		nodes[i] = scheduler.NewNode(i, nc.Url)
	}
	return nodes
}
