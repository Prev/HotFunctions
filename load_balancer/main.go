package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	nodes := initNodesFromConfig("nodes.config.json")

	fileEvents, err := os.Open("data/events.csv")
	if err != nil {
		panic(err)
	}
	defer fileEvents.Close()
	rdr := csv.NewReader(bufio.NewReader(fileEvents))
	rows, _ := rdr.ReadAll()

	os.MkdirAll("logs", 0755)
	logFileName := "logs/log-" + time.Now().Format("2006-01-02T15:04:05") + ".out"
	outputFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	logger = log.New(outputFile, "", log.Ldate|log.Ltime)

	sched := newLeastLoadedScheduler(&nodes)
	// sched := newConsistentHashingScheduler(&nodes, 8)
	// sched := newTableBasedScheduler(&nodes)
	// sched := newAdaptiveScheduler(&nodes, 2)

	stdTick := time.Now().UnixNano() / int64(time.Millisecond)

	tick := 0
	i := 0
	for tick <= 60000 {
		if i >= len(rows) {
			break
		}
		row := rows[i]
		startTime, _ := strconv.Atoi(row[1])

		tick = int((time.Now().UnixNano() / int64(time.Millisecond)) - stdTick)

		if tick >= startTime {
			i += 1
			name := row[0]
			node, err := sched.pick(name)

			if err != nil {
				fmt.Printf("[fail running] %s: %s\n", name, err.Error())
				continue
			}

			go node.runFunction(name, func(n string, time int64) {
				// sched.appendExecutionResult(n, time)
			})

		} else {
			time.Sleep(time.Second / 100)
		}
	}

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
