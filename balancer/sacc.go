package balancer

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
	"strconv"
)

type SACC struct {
	Context ExpContext
	WithCSTx bool
	CollationUtils [][]BalanceInfo
}

func FindMinMax(slice []BalanceInfo) (int64, int64) {
	min := int64(10000000)
	max := int64(0)
	for i := 0; i < len(slice); i++ {
		if min > slice[i].GasUsed {
			min = slice[i].GasUsed
		}

		if max < slice[i].GasUsed {
			max = slice[i].GasUsed
		}
	}

	return min, max
}

func (s *SACC) Init(withCSTx bool) {
	s.WithCSTx = withCSTx
	s.CollationUtils = make([][]BalanceInfo, s.Context.CollationCycle)
	for i := 0; i < s.Context.CollationCycle; i++ {
		s.CollationUtils[i] = make([]BalanceInfo, s.Context.NumberOfShards)
	}
}

func (s *SACC) GetUtilization() float64 {
	var utilization float64

	for i := 0; i< s.Context.CollationCycle; i++ {
		var sum int64
		for j := 0; j < s.Context.NumberOfShards; j++ {
			sum += s.CollationUtils[i][j].GasUsed
		}

		currentUtilization := float64(sum) / float64(s.Context.GasLimit * s.Context.NumberOfShards)
		utilization += (currentUtilization - utilization) / float64(i+1)
	}

	return utilization
}

func (s *SACC) GetMakespan() float64 {
	var makespan float64
	for i := 0; i < s.Context.CollationCycle; i++ {
		min, max := FindMinMax(s.CollationUtils[i])
		makespan += (float64(max-min)-makespan) / float64(i+1)
	}

	return makespan
}

func (s *SACC) GetThroughput() float64 {
	var avg float64

	for i := 0; i < s.Context.CollationCycle; i++ {
		var sum float64
		for j := 0; j < s.Context.NumberOfShards; j++ {
			sum += float64(s.CollationUtils[i][j].Transactions)
		}

		avg += (sum - avg) / float64(i+1)
	}

	return avg
}

func (s *SACC) GetCrossShards() float64 {
	var sum float64

	for i := 0; i < s.Context.CollationCycle; i++ {
		for j := 0; j < s.Context.NumberOfShards; j++ {
			sum += float64(s.CollationUtils[i][j].CrossShards)
		}
	}

	return sum/float64(s.Context.CollationCycle)
}

func (s *SACC) StartExperiment() {
	db := storage.OpenDatabase()
	var transaction storage.Transaction

	testNumber := 0
	toBlock := s.Context.FromBlock + (s.Context.CollationCycle * s.Context.BlockEpoch)
	for i := s.Context.FromBlock; i < toBlock; i++ {
		cursor := db.FetchTransactions(int64(i))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&transaction)
			utils.CheckError(err, "failed to decode transaction data")

			shardNum := utils.GetShardSaccAddress(transaction.ToAddress, s.Context.NumberOfShards)
			senderShard := utils.GetShardSaccAddress(transaction.Sender, s.Context.NumberOfShards)

			// select shard.
			sd := &s.CollationUtils[testNumber][shardNum]
			conflictSd := &s.CollationUtils[testNumber][senderShard]
			limited := sd.GasUsed + transaction.GasUsed < int64(s.Context.GasLimit)

			if s.WithCSTx && testNumber > 0 {
				prevSd := s.CollationUtils[testNumber-1][shardNum]
				limited = (sd.GasUsed + transaction.GasUsed + int64(prevSd.CrossShards * 42000)) < int64(s.Context.GasLimit)
			}

			if limited {
				sd.GasUsed += transaction.GasUsed
				sd.ElapsedTime += transaction.ElapsedTime
				sd.Transactions += 1
				if shardNum != senderShard {
					conflictSd.CrossShards += 1
				}
			}
		}

		if (i+1) % s.Context.BlockEpoch == 0 {
			testNumber += 1
		}
	}

	fmt.Printf("[S-ACC] Collation Utilization : %.3f%%\n", s.GetUtilization() * 100)
	fmt.Printf("[S-ACC] Normalized Makespan: %.3f\n", s.GetMakespan())
	fmt.Printf("[S-ACC] Average Transaction Throuput: %.3f\n", s.GetThroughput())
	fmt.Printf("[S-ACC] Average Cross-shard txes: %.3f\n", s.GetCrossShards())
}

func (s *SACC) SaveToCSV(filename string){
	data := make([][]string, s.Context.CollationCycle + 1)
	for i := 0; i < s.Context.CollationCycle + 1; i++ {
		data[i] = make([]string, s.Context.NumberOfShards)
	}

	for j := 0; j < s.Context.NumberOfShards; j++ {
		data[0][j] = "shard " + strconv.Itoa(j)
	}

	for i := 1; i <= s.Context.CollationCycle; i++ {
		for j := 0; j < s.Context.NumberOfShards; j++ {
			data[i][j] = strconv.FormatInt(s.CollationUtils[i-1][j].GasUsed, 10)
		}
	}

	utils.ReportToCSV(filename, data)
}