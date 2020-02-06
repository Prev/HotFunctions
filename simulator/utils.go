package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

func startSimulation(eventStreamPath string, action func(string, int)) {
	file, err := os.Open(eventStreamPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(bufio.NewReader(file))
	rows, err := rdr.ReadAll()
	if err != nil {
		panic(err)
	}

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
			go action(row[0], startTime)

		} else {
			time.Sleep(time.Second / 100)
		}
	}
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
