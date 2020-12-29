package balancer

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
)

type SACC struct {
	FromBlock int
	ToBlock int
	BlockEpoch int
	NumberOfShards int
	GasLimit int
	CollationUtils [][]int64
}

func (s *SACC) Init() {
	totalTests := (s.ToBlock - s.FromBlock)/20
	s.CollationUtils = make([][]int64, totalTests)
	for i := 0; i < totalTests; i++ {
		s.CollationUtils[i] = make([]int64, s.NumberOfShards)
	}
}

func (s *SACC) StartExperiment() {
	db := storage.OpenDatabase()
	var transaction storage.Transaction

	testNumber := 0
	for i := s.FromBlock; i < s.ToBlock; i++ {
		cursor := db.FetchTransactions(int64(i))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&transaction)
			utils.CheckError(err, "failed to decode transaction data")

			shardNum := utils.GetShardSaccAddress(transaction.ToAddress, s.NumberOfShards)
			// if gas limit hit.
			if s.CollationUtils[testNumber][shardNum] < int64(s.GasLimit) {
				s.CollationUtils[testNumber][shardNum] += transaction.GasUsed
			}
		}

		if (i+1) % s.BlockEpoch == 0 {
			fmt.Println(s.CollationUtils[testNumber])
			testNumber += 1
		}
	}
}