package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Prev/LALB/load_balancer/scheduler"
)

type NodeConfigData struct {
	Url         string `json:"url"`
	MaxCapacity int    `json:"maxCapacity"`
}

var logger *log.Logger
var mutex *sync.Mutex

func main() {
	if len(os.Args) != 2 {
		println("Usage: go run *.go ll|ch|tb|ad")
		os.Exit(-1)
	}

	nodes := initNodesFromConfig("nodes.config.json")

	var sched scheduler.Scheduler
	switch os.Args[1] {
	case "ll":
		println("Using Least Loaded Scheduler")
		sched = scheduler.NewLeastLoadedScheduler(&nodes)
	case "ch":
		println("Using Consistent Hashing Scheduler")
		sched = scheduler.NewConsistentHashingScheduler(&nodes, 8)
	case "tb":
		println("Using Table-based Scheduler")
		sched = scheduler.NewTableBasedScheduler(&nodes)
	case "ad":
		println("Using Adaptive Scheduler")
		sched = scheduler.NewAdaptiveScheduler(&nodes, 2)
	}

	logger = initLogger()

	simulator := newSimulator("data/events.csv")
	simulator.Start(func(functionName string, time int) {
		if functionName != "W3" {
			return
		}
		node, err := sched.Select(functionName)

		if err != nil {
			fmt.Printf("[fail running] %s: %s\n", functionName, err.Error())
			return
		}

		go runFunction(node, functionName, func(name string, time int64) {
			sched.Finished(node, name, time)
		})
	})

	time.Sleep(time.Second * 20)
}

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
		nodes[i] = scheduler.NewNode(i, nc.Url, nc.MaxCapacity)
	}
	return nodes
}

func initLogger() *log.Logger {
	dirName := "logs/" + time.Now().Format("2006-01-02")
	os.MkdirAll(dirName, 0755)

	logFileName := dirName + "/" + time.Now().Format("15:04:05") + ".log"
	outputFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return log.New(outputFile, "", log.Ldate|log.Ltime)
}
