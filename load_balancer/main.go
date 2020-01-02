package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type NodeConfigData struct {
	Url         string `json:"url"`
	MaxCapacity int    `json:"maxCapacity"`
}

var logger *log.Logger
var mutex *sync.Mutex

func main() {
	logger = initLogger()
	nodes := initNodesFromConfig("nodes.config.json")
	sched := newLeastLoadedScheduler(&nodes)
	// sched := newConsistentHashingScheduler(&nodes, 8)
	// sched := newTableBasedScheduler(&nodes)
	// sched := newAdaptiveScheduler(&nodes, 2)

	simulator := newSimulator("data/events.csv")
	simulator.Start(func(functionName string, time int) {
		node, err := sched.pick(functionName)

		if err != nil {
			fmt.Printf("[fail running] %s: %s\n", functionName, err.Error())
			return
		}

		go node.runFunction(functionName, func(n string, time int64) {
			// sched.appendExecutionResult(n, time)
		})
	})

	time.Sleep(time.Second * 20)
}

func initNodesFromConfig(configFilePath string) []*Node {
	nodeConfigFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer nodeConfigFile.Close()

	var nodeConfigs []NodeConfigData
	byteValue, _ := ioutil.ReadAll(nodeConfigFile)
	json.Unmarshal([]byte(byteValue), &nodeConfigs)

	nodes := make([]*Node, len(nodeConfigs))
	for i, nc := range nodeConfigs {
		nodes[i] = newNode(i, nc.Url, nc.MaxCapacity)
	}
	return nodes
}

func initLogger() *log.Logger {
	os.MkdirAll("logs", 0755)
	logFileName := "logs/log-" + time.Now().Format("2006-01-02T15:04:05") + ".out"
	outputFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return log.New(outputFile, "", log.Ldate|log.Ltime)
}

func printCapacityTable(nodes *[]*Node) {
	out := "--------capacity table--------\n" +
		"|Node\t|Total\t|Spare\t|Run\t|\n"

	for _, node := range *nodes {
		runningFunctions := ""
		for fname, num := range node.running {
			for i := 0; i < num; i++ {
				runningFunctions += fmt.Sprintf("%s, ", fname)
			}
		}

		out += fmt.Sprintf("|%d\t|%d\t|%d\t|%s\t|\n", node.id, node.maxCapacity, node.capacity(), runningFunctions)
	}
	out += "--------------------------"
	println(out)
}
