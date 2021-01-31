package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/Prev/HotFunctions/load_balancer/scheduler"
)

var logger *log.Logger
var sched scheduler.Scheduler
var schedType string
var fakeMode = false
var nodes []*scheduler.Node

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 2 {
		println("Usage: go run *.go rr|ll|hash|ours|pasch|tradeoff [fakeMode=0]")
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
	setScheduler(schedType)

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

func setScheduler(newSchedType string) error {
	var err error

	switch newSchedType {
	case "rr":
		// Round Robin Scheduler
		println("Using Round Robin Scheduler")
		sched = scheduler.NewRoundRobinScheduler(&nodes)
	case "ll":
		// Least Loaded Scheduler
		// Scheduler picks the node who has minimum executing tasks
		println("Using Least Loaded Scheduler")
		sched = scheduler.NewLeastLoadedScheduler(&nodes)

	case "hash":
		// Consistent Hashing Scheduler
		// Scheduler picks the node by Consistent Hashing algorithm where key is the function name
		var threshold int
		if threshold, err = strconv.Atoi(os.Getenv("T")); err != nil {
			threshold = 8
		}
		fmt.Printf("Using Consistent Hashing Scheduler (t=%d)\n", threshold)
		sched = scheduler.NewConsistentHashingScheduler(&nodes, threshold, uint(threshold))

	case "pasch":
		// Consistent Hashing Scheduler
		// Scheduler picks the node by Consistent Hashing algorithm where key is the function name
		var threshold int
		if threshold, err = strconv.Atoi(os.Getenv("T")); err != nil {
			threshold = 8
		}
		fmt.Printf("Using PASch Extended Scheduler (t=%d)\n", threshold)
		sched = scheduler.NewPASchExtendedScheduler(&nodes, uint(threshold))

	case "ours":
		// Proposing Greedy Scheduler
		var t1, t2, t3 int
		if t1, err = strconv.Atoi(os.Getenv("T1")); err != nil {
			t1 = 8
		}
		if t2, err = strconv.Atoi(os.Getenv("T2")); err != nil {
			t2 = 5
		}
		if t3, err = strconv.Atoi(os.Getenv("T3")); err != nil {
			t3 = 3
		}
		println("Using Our Scheduler")
		sched = scheduler.NewOurScheduler(&nodes, uint(t1), uint(t2), t3)

	case "tradeoff":
		// Trade-off algorithm for exploiting locality
		var alpha, beta float64
		if alpha, err = strconv.ParseFloat(os.Getenv("ALPHA"), 32); err != nil {
			alpha = 0.0
		}
		if beta, err = strconv.ParseFloat(os.Getenv("BETA"), 32); err != nil {
			beta = 0.5
		}
		fmt.Printf("Using Trade-Off Scheduler (alpha=%.3f, beta=%.3f)\n", alpha, beta)
		sched = scheduler.NewTradeOffScheduler(&nodes, alpha, beta)

	default:
		return errors.New("unsupported scheduler type")
	}

	schedType = newSchedType
	return nil
}
