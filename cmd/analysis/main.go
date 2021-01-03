package main

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
	"strconv"
)

func main() {
	runAccountWorkload("analysis_account_workload.csv")
	runTransactionWorkload("analysis_transaction_workload.csv")
}

// save the amount of transactions that occurred for each account group.
func runAccountWorkload(filename string) {
	const AccountGroup = 100
	const FromBlock = 7000000
	const ToBlock = 7000050

	var tx storage.Transaction

	txList := make([]int64, AccountGroup)
	gasList := make([]int64, AccountGroup)

	db := storage.OpenDatabase()
	for i := FromBlock; i <= ToBlock; i++ {
		cursor := db.FetchTransactions(int64(i))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&tx)
			utils.CheckError(err, "failed to decode transaction")

			num := utils.GetShardSaccAddress(tx.ToAddress, AccountGroup)
			txList[num]++
			gasList[num] += int64(tx.GasUsed)
		}
	}

	data := make([][]string, AccountGroup)

	for i := 0; i < AccountGroup; i++ {
		data[i] = make([]string, 2)
		data[i][0] = strconv.Itoa(i)
		data[i][1] = strconv.FormatInt(txList[i], 10)
	}

	utils.ReportToCSV(filename, data)
	fmt.Printf("-- saved %s.\n", filename)
}

func runTransactionWorkload(filename string) {

}