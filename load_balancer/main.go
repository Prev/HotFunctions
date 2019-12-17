package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	logger := newLogger(&nodes)

	file, err := os.Open("data/events.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()

	i := 0
	for tick := 0; tick <= 60000; {
		logger.logDistinctFunctionsCounts()

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
			} else if name == "W1" || name == "W2" || name == "W3" || name == "W4" || name == "W5" {
				go node.runFunction(name)
			}
			i += 1

		} else {
			tick += 10
			time.Sleep(time.Second / 3)
		}
	}

	for i, _ := range nodes {
		sum := 0
		for j := 0; j < logger.curIndex; j++ {
			sum += logger.distinctFuncs[i][j]
		}
		fmt.Printf("node %d: %.2f\n", i, float32(sum)/float32(logger.curIndex))
	}
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

	tmp := make([]*Node, len(nodeConfigs))
	for i, nc := range nodeConfigs {
		tmp[i] = newNode(i, nc.Url, nc.MaxCapacity)
	}
	// nodes = &tmp
	return tmp
}
