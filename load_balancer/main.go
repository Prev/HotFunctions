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
	"time"
)

type NodeConfigData struct {
	Url         string `json:"url"`
	MaxCapacity int    `json:"maxCapacity"`
}

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
	myLogger := log.New(outputFile, "", log.Ldate|log.Ltime)

	num := 0
	i := 0
	for tick := 0; tick <= 20000; {
		if i >= len(rows) {
			break
		}
		row := rows[i]
		startTime, _ := strconv.Atoi(row[1])

		if tick >= startTime {
			name := row[0]
			node, err := leastLoaded(&nodes)

			if err != nil {
				fmt.Printf("[fail running] %s: %s\n", name, err.Error())
			} else if name == "W1" || name == "W2" || name == "W3" || name == "W4" || name == "W5" || name == "W6" || name == "W7" {
				num += 1
				go node.runFunction(name, myLogger)
			}
			i += 1

		} else {
			tick += 10
			time.Sleep(time.Second / 100)
		}
	}

	time.Sleep(time.Second * 20)
	println(num)
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
