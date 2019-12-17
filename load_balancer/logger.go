package main

type Logger struct {
	distinctFuncs [10][1000]int
	curIndex      int
	nodes         *[]*Node
}

func newLogger(nodes *[]*Node) *Logger {
	logger := new(Logger)
	logger.curIndex = 0
	logger.nodes = nodes
	return logger
}

func (logger *Logger) logDistinctFunctionsCounts() {
	logger.curIndex += 1
	for i, n := range *logger.nodes {
		distincts := make(map[string]bool)
		for _, f := range n.running {
			distincts[f] = true
		}
		logger.distinctFuncs[i][logger.curIndex] = len(distincts)
	}
}
