package main

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
)

var (
	FromBlock = 6000000
	ToBlock = 7000000
	BlockEpoch = 20
	NumberOfShards = 20
	GasLimit = 10000000
)

func GetUtilization() {

}

func GetMakespan() {
	
}

func main() {
	fmt.Println("Static Address Allocation (S-ACC) experiment.")

	// 테스트 데이터를 저장할 배열 생성
	totalTests := (ToBlock - FromBlock)/20
	collationUtils := make([][]int64, totalTests)
	for i := 0; i < totalTests; i++ {
		collationUtils[i] = make([]int64, NumberOfShards)
	}

	db := storage.OpenDatabase()
	var transaction storage.Transaction

	testNumber := 0
	for i := FromBlock; i < ToBlock; i++ {
		cursor := db.FetchTransactions(int64(i))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&transaction)
			utils.CheckError(err, "failed to decode transaction data")

			shardNum := utils.GetShardSaccAddress(transaction.ToAddress, NumberOfShards)
			// if gas limit hit.
			if collationUtils[testNumber][shardNum] < int64(GasLimit) {
				collationUtils[testNumber][shardNum] += transaction.GasUsed
			}
		}

		if (i+1) % BlockEpoch == 0 {
			fmt.Println(collationUtils[testNumber])
			testNumber += 1
		}
	}
}