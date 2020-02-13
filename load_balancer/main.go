package main

import (
	"encoding/json"
	"fmt"
	"github.com/Prev/HotFunctions/load_balancer/scheduler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var logger *log.Logger
var sched scheduler.Scheduler
var schedType string
var fakeMode = false
var nodes []*scheduler.Node

func main() {
	if len(os.Args) < 2 {
		println("Usage: go run *.go ll|hash|ours|pasch [fakeMode=0]")
		os.Exit(-1)
	}
	schedType = os.Args[1]

	if len(os.Args) > 2 {
		intVal, err := strconv.Atoi(os.Args[2])
		if err != nil {
			println("Argument fakeMode should be integer (its value means number of the fake nodes)")
			os.Exit(-1)
		}

		if intVal > 0 {
			println("Run as fakeMode (do not send a request to the worker node), # of nodes: ", intVal)
			fakeMode = true
			for i := 0; i < intVal; i++ {
				nodes = append(nodes, scheduler.NewNode(i, ""))
			}
			goto GuessSchedType
		}
	}
	println("Load node info from `nodes.config.json`")
	nodes = initNodesFromConfig("nodes.config.json")
	fmt.Printf("%d nodes found\n", len(nodes))

GuessSchedType:
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
		sched = scheduler.NewConsistentHashingScheduler(&nodes, 8, 8)

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

	var nodeConfigs []string
	byteValue, _ := ioutil.ReadAll(nodeConfigFile)
	json.Unmarshal(byteValue, &nodeConfigs)

	nodes := make([]*scheduler.Node, len(nodeConfigs))
	for i, url := range nodeConfigs {
		nodes[i] = scheduler.NewNode(i, url)
	}
	return nodes
}
