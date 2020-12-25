package reporter

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
	"os"
	"sync"
	"time"
)

const (
	StartBlock = 5000000
	SampleBlockRate = 25
	PrintEpoch = 10000000
	PipelineSize = 1000000
	NumberOfWorkers = 50
)

func ReportToCSV(acc map[string]float64, cnt map[string]int64) {
	file, err := os.Create("opcodeAvg.csv")
	utils.CheckError(err, "failed to open opcodeAvg.csv file")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for key, value := range acc {
		data := []string { key, fmt.Sprintf("%f", value), fmt.Sprintf("%d", cnt[key])}
		err := writer.Write(data)
		utils.CheckError(err, "failed to write opcodeAvg.csv file")
	}
}

func ReportAvgStd(done chan bool) {
	fmt.Println("[AvgStdReporter] started")

	accChannel := make(chan map[string]float64, 1)
	cntChannel := make(chan map[string]int64, 1)
	var wg sync.WaitGroup

	pipeline := make(chan storage.Opcode, PipelineSize)
	db := storage.OpenDatabase()

	for i := 1; i <= NumberOfWorkers; i++ {
		wg.Add(1)
		go doWork(i, db.MaxBlockNumber(), pipeline, &wg)
	}
	go accumulator(pipeline, accChannel, cntChannel)
	wg.Wait()
	_acc := <-accChannel
	_cnt := <-cntChannel
	ReportToCSV(_acc, _cnt)
	done <- true
}

// worker
func doWork(workerId int, maxBlockNumber int64, pipeline chan<- storage.Opcode, wg *sync.WaitGroup) {
	defer wg.Done()
	var opcode storage.Opcode

	fmt.Printf("worker %d ready\n", workerId)
	blockNumber := int64(workerId) * SampleBlockRate + StartBlock

	for blockNumber <= maxBlockNumber {
		db := storage.OpenDatabase()
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
func accumulator(pipeline <-chan storage.Opcode, accChannel chan map[string]float64, cntChannel chan map[string]int64) {
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
			for k,v := range opcodeAvg {
				fmt.Printf("[%s]:%f, count : %d", k, v, opcodeCount[k])
			}
		}
	}

	accChannel <- opcodeAvg
	cntChannel <- opcodeCount
}
