package main

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
)

func main() {
	var tx storage.Transaction

	AccountGroup := 20
	FromBlock := 7000000
	ToBlock := 7000020
	txList := make([]uint64, AccountGroup)
	gasList := make([]uint64, AccountGroup)

	db := storage.OpenDatabase()
	for i := FromBlock; i <= ToBlock; i++ {
		cursor := db.FetchTransactions(int64(i))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&tx)
			utils.CheckError(err, "failed to decode transaction")

			num := utils.GetShardSaccAddress(tx.ToAddress, AccountGroup)
			txList[num]++
			gasList[num] += uint64(tx.GasUsed)
		}
	}

	for i := 0; i < AccountGroup; i++{
		fmt.Printf("[%d] tx: %d, gas: %d\n", i, txList[i], gasList[i])
	}
}
