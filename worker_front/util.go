package main

import (
	"os"
	"sort"
	"strconv"
)

type pair struct {
	value int64
	key string
}

func sortMapByValue(mapObject *map[string]int64) []string {
	n := len(*mapObject)
	tmp := make([]pair, 0, n)

	for key, val := range *mapObject {
		tmp = append(tmp, pair{val, key})
	}
	sort.Slice(tmp, func (i, j int) bool {
		return tmp[i].value > tmp[j].value
	})

	ret := make([]string, n)
	for i, item := range tmp {
		ret[i] = item.key
	}
	return ret
}


func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}