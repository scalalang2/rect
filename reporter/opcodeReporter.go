package reporter

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"time"
)

const (
	StartBlock = 5000000
	SampleBlockRate = 25
	PrintEpoch = 100000
	PipelineSize = 1000000
	NumberOfWorkers = 10
)

func ReportAvgStd(done chan bool) {
	fmt.Println("[AvgStdReporter] started")
	pipeline := make(chan storage.Opcode, PipelineSize)
	db := storage.OpenDatabase()

	for i := 1; i <= NumberOfWorkers; i++ {
		go doWork(i, db.MaxBlockNumber(), pipeline)
	}
	go accumulator(pipeline)
}

// worker
func doWork(workerId int, maxBlockNumber int64, pipeline chan<- storage.Opcode) {
	var opcode storage.Opcode

	fmt.Printf("worker %d ready\n", workerId)
	db := storage.OpenDatabase()
	blockNumber := int64(workerId) * SampleBlockRate + StartBlock

	for blockNumber <= maxBlockNumber {
		cursor := db.FindOpcodes(blockNumber)
		for cursor.Next(context.TODO()) {
			cursor.Decode(&opcode)
			opcode.WorkerId = workerId
			pipeline <- opcode
		}
		blockNumber = blockNumber + NumberOfWorkers * SampleBlockRate
	}
}

// accumulator
// formula: avg += (x - avg) / n;
func accumulator(pipeline <-chan storage.Opcode) {
	opcodeAvg := make(map[string]float64, 100)
	opcodeCount := make(map[string]int64, 100)
	receivedWorks := make(map[int]int64, NumberOfWorkers)

	startTime := time.Now().Unix()

	i := 0
	for opcode := range pipeline {
		_, doesExist := opcodeAvg[opcode.Opcode]
		if !doesExist {
			opcodeAvg[opcode.Opcode] = 0.0
			opcodeCount[opcode.Opcode] = 0
		}

		_, doesExist = receivedWorks[opcode.WorkerId]
		if !doesExist {
			receivedWorks[opcode.WorkerId] = 0
		}

		opcodeCount[opcode.Opcode]++
		opcodeAvg[opcode.Opcode] += (float64(opcode.ElapsedTime) - opcodeAvg[opcode.Opcode]) / float64(opcodeCount[opcode.Opcode])
		receivedWorks[opcode.WorkerId]++

		i++
		if (i+1) % PrintEpoch == 0 {
			endTime := time.Now().Unix()
			fmt.Printf("records: %d, acc elapsed time: [%ds]\n", i+1, endTime - startTime)
		}
	}
}
