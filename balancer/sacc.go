package balancer

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
	"strconv"
)

type BalanceInfo struct {
	GasUsed int64
	ElapsedTime int64
	Transactions int
	CrossShards int
}

type SACC struct {
	FromBlock int
	ToBlock int
	BlockEpoch int
	NumberOfShards int
	GasLimit int
	TotalTests int
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
	s.TotalTests = (s.ToBlock - s.FromBlock)/s.NumberOfShards
	s.CollationUtils = make([][]BalanceInfo, s.TotalTests)
	for i := 0; i < s.TotalTests; i++ {
		s.CollationUtils[i] = make([]BalanceInfo, s.NumberOfShards)
	}
}

func (s *SACC) GetUtilization() float64 {
	var utilization float64

	for i := 0; i< s.TotalTests; i++ {
		var sum int64
		for j := 0; j < s.NumberOfShards; j++ {
			sum += s.CollationUtils[i][j].GasUsed
		}

		currentUtilization := float64(sum) / float64(s.GasLimit * s.NumberOfShards)
		utilization += (currentUtilization - utilization) / float64(i+1)
	}

	return utilization
}

func (s *SACC) GetMakespan() float64 {
	var makespan float64
	for i := 0; i < s.TotalTests; i++ {
		min, max := FindMinMax(s.CollationUtils[i])
		makespan += (float64(max-min)-makespan) / float64(i+1)
	}

	return makespan
}

func (s *SACC) GetThroughput() float64 {
	var avg float64

	for i := 0; i < s.TotalTests; i++ {
		var sum float64
		for j := 0; j < s.NumberOfShards; j++ {
			sum += float64(s.CollationUtils[i][j].Transactions)
		}

		avg += (sum - avg) / float64(i+1)
	}

	return avg
}

func (s *SACC) GetCrossShards() float64 {
	total := (s.ToBlock - s.FromBlock)/20
	var sum float64

	for i := 0; i < total; i++ {
		for j := 0; j < s.NumberOfShards; j++ {
			sum += float64(s.CollationUtils[i][j].CrossShards)
		}
	}

	return sum/float64(total)
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
			senderShard := utils.GetShardSaccAddress(transaction.Sender, s.NumberOfShards)

			// select shard.
			sd := &s.CollationUtils[testNumber][shardNum]
			conflictSd := &s.CollationUtils[testNumber][senderShard]
			limited := sd.GasUsed + transaction.GasUsed < int64(s.GasLimit)

			if s.WithCSTx && testNumber > 0 {
				prevSd := s.CollationUtils[testNumber-1][shardNum]
				limited = (sd.GasUsed + transaction.GasUsed + int64(prevSd.CrossShards * 42000)) < int64(s.GasLimit)
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

		if (i+1) % s.BlockEpoch == 0 {
			testNumber += 1
		}
	}

	fmt.Printf("[S-ACC] Collation Utilization : %.3f%%\n", s.GetUtilization() * 100)
	fmt.Printf("[S-ACC] Normalized Makespan: %.3f\n", s.GetMakespan())
	fmt.Printf("[S-ACC] Average Transaction Throuput: %.3f\n", s.GetThroughput())
	fmt.Printf("[S-ACC] Average Cross-shard txes: %.3f\n", s.GetCrossShards())
}

func (s *SACC) SaveToCSV(filename string){
	data := make([][]string, s.TotalTests + 1)
	for i := 0; i < s.TotalTests + 1; i++ {
		data[i] = make([]string, s.NumberOfShards)
	}

	for j := 0; j < s.NumberOfShards; j++ {
		data[0][j] = "shard " + strconv.Itoa(j)
	}

	for i := 1; i <= s.TotalTests; i++ {
		for j := 0; j < s.NumberOfShards; j++ {
			data[i][j] = strconv.FormatInt(s.CollationUtils[i-1][j].GasUsed, 10)
		}
	}

	utils.ReportToCSV(filename, data)
}