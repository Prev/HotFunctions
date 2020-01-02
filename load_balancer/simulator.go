package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type Simulator struct {
	eventStreamPath string
}

func newSimulator(eventStreamPath string) *Simulator {
	sim := new(Simulator)
	sim.eventStreamPath = eventStreamPath
	return sim
}

func (s *Simulator) Start(action func(string, int)) {
	file, err := os.Open(s.eventStreamPath)
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
			action(row[0], startTime)

		} else {
			time.Sleep(time.Second / 100)
		}
	}
}
