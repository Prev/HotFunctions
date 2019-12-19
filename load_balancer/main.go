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

	// sched := newTableBasedScheduler()
	sched := newAdaptiveScheduler(2)

	i := 0
	for tick := 0; tick <= 20000; {
		if i >= len(rows) {
			break
		}
		row := rows[i]
		startTime, _ := strconv.Atoi(row[1])

		if tick >= startTime {
			i += 1
			name := row[0]

			if name != "W1" && name != "W2" && name != "W3" && name != "W4" && name != "W5" && name != "W6" && name != "W7" {
				continue
			}

			// node, err := leastLoaded(&nodes)
			node, err := sched.pick(&nodes, name)

			if err != nil {
				// fmt.Printf("[fail running] %s: %s\n", name, err.Error())
				continue
			}

			go node.runFunction(name, func(n string, time int64) {
				sched.appendExecutionResult(n, time)
			})

		} else {
			tick += 10
			time.Sleep(time.Second / 100)
		}

		if tick%1000 == 0 {
			sched.printTables()
			printCapacityTable(&nodes)
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
