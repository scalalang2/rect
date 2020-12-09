package reporter

import (
	"context"
	"fmt"
	"load-balancing-simulator/storage"
	"time"
)

func ReportAvgStd(done chan bool) {
	fmt.Println("[AvgStdReporter] started")
	numberOfWorkers := 10
	pipeline := make(chan storage.Opcode, 100000)
	db := storage.OpenDatabase()

	for i := 0; i < numberOfWorkers; i++ {
		go doWork(i, numberOfWorkers, db.MaxBlockNumber(), pipeline)
	}
	go accumulator(pipeline)
}

// worker
func doWork(workerId int, workersNum int, maxBlockNumber int64, pipeline chan<- storage.Opcode) {
	var opcode storage.Opcode

	fmt.Printf("worker %d ready\n", workerId)
	db := storage.OpenDatabase()
	blockNumber := int64(workerId)

	for blockNumber <= maxBlockNumber {
		cursor := db.FindOpcodes(blockNumber)
		for cursor.Next(context.TODO()) {
			cursor.Decode(&opcode)
			opcode.WorkerId = workerId
			pipeline <- opcode
		}
		blockNumber = blockNumber + int64(workersNum)
	}
}

// accumulator
// formula: avg += (x - avg) / n;
func accumulator(pipeline <-chan storage.Opcode) {
	opcodeAvg := make(map[string]float64, 50)
	opcodeCount := make(map[string]int64, 50)
	receivedWorks := make(map[int]int64, 10)

	startTime := time.Now().Unix()

	i := 0
	for opcode := range pipeline {
		_, doesExist := opcodeAvg[opcode.Opcode]
		if !doesExist {
			opcodeAvg[opcode.Opcode] = 0.0
			opcodeCount[opcode.Opcode] = 0
		}
		opcodeCount[opcode.Opcode]++
		opcodeAvg[opcode.Opcode] += (float64(opcode.ElapsedTime) - opcodeAvg[opcode.Opcode]) / float64(opcodeCount[opcode.Opcode])
		receivedWorks[opcode.WorkerId]++

		i++
		if (i+1) % 1000 == 0 {
			endTime := time.Now().Unix()
			fmt.Printf("records: %d, acc elapsed time: [%ds]\n", i+1, endTime - startTime)
			for k, v := range opcodeAvg {
				fmt.Printf("%s : %.6f\n", k, v)
			}

			for k, v := range receivedWorks {
				fmt.Printf("workerId: %d, jobs: %d\n", k, v)
			}
		}
	}
}
