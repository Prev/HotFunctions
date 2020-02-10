package main

import (
	"os"
	"sort"
	"strconv"
)

type lruPair struct {
	value int64
	key string
}

func leastRecentlyUsed(lru *map[string]int64, n int) ([]string, []string) {
	numTotalItems := len(*lru)
	if numTotalItems < n {
		n = numTotalItems
	}

	tmp := make([]lruPair, 0, numTotalItems)

	for key, val := range *lru {
		tmp = append(tmp, lruPair{val, key})
	}

	sort.Slice(tmp, func (i, j int) bool {
		return tmp[i].value > tmp[j].value
	})

	live := make([]string, 0, n)
	dead := make([]string, 0, numTotalItems - n)

	for i, item := range tmp {
		if i < n {
			live = append(live, item.key)
		} else {
			dead = append(dead, item.key)
		}
	}

	return live, dead
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