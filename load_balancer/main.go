package main

import (
	"encoding/json"
	"fmt"
	"github.com/Prev/HotFunctions/load_balancer/scheduler"
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

	nodes := initNodesFromConfig("nodes.config.json")

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
