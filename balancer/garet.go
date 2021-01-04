package balancer

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
)

type GARET struct {
	Context ExpContext
	AccountGroups int
	Ncol int
	WithCSTx bool
	CollationUtils [][]BalanceInfo
	GasUsedAcc [][]int
	GasPredAcc []int
	GasPredShard []int
	MappingTable map[int]int // key: i-th account group. value: j-th shard
}

func (g *GARET) Init(withCSTx bool) {
	g.WithCSTx = withCSTx
	g.CollationUtils = make([][]BalanceInfo, g.Context.CollationCycle)
	g.MappingTable = make(map[int]int, g.AccountGroups)

	// amount of gas used for each account group.
	g.GasUsedAcc = make([][]int, g.AccountGroups)
	g.GasPredAcc = make([]int, g.AccountGroups)
	g.GasPredShard = make([]int, g.Context.NumberOfShards)

	for i := 0; i < g.Context.CollationCycle; i++ {
		g.CollationUtils[i] = make([]BalanceInfo, g.Context.NumberOfShards)
	}

	for i := 0; i < g.AccountGroups; i++ {
		g.GasUsedAcc[i] = make([]int, g.Ncol)
	}

	for i := 0; i < g.AccountGroups; i++ {
		g.MappingTable[i] = i % g.Context.NumberOfShards
	}
}

func (g *GARET) StartExperiment() {
	db := storage.OpenDatabase()
	var tx storage.Transaction

	utilNumber := 0
	toBlock := g.Context.FromBlock + (g.Context.CollationCycle * g.Context.BlockEpoch)
	for blockNumber := g.Context.FromBlock; blockNumber < toBlock; blockNumber++ {
		cursor := db.FetchTransactions(int64(blockNumber))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&tx)
			utils.CheckError(err, "failed to decode transactions")

			toAccGroup := utils.GetShardSaccAddress(tx.ToAddress, g.AccountGroups)
			fromAccGroup := utils.GetShardSaccAddress(tx.Sender, g.AccountGroups)
			toShardNum := g.MappingTable[toAccGroup]
			fromShardNum := g.MappingTable[fromAccGroup]

			// util 번 째 블록의 to, from 샤드
			toShard := &g.CollationUtils[utilNumber][toShardNum]
			fromShard := &g.CollationUtils[utilNumber][fromShardNum]

			notLimited := toShard.GasUsed + tx.GasUsed < int64(g.Context.GasLimit)

			// consider cross-shard tx
			if g.WithCSTx && utilNumber > 0 {
				prevShard := g.CollationUtils[utilNumber-1][toShardNum]
				notLimited = toShard.GasUsed + tx.GasUsed + int64(prevShard.CrossShards * 42000) < int64(g.Context.GasLimit)
			}

			if notLimited {
				g.GasUsedAcc[toAccGroup][utilNumber % g.Ncol] += int(tx.GasUsed)
				toShard.GasUsed += tx.GasUsed
				toShard.Transactions += 1
				if toShardNum != fromShardNum {
					fromShard.CrossShards += 1
				}
			}
		}

		if (blockNumber + 1) % g.Context.BlockEpoch == 0 {
			utilNumber += 1
			if utilNumber % g.Ncol == 0 {
				g.AccGroupRelocation()
			}
			fmt.Println(g.CollationUtils[utilNumber-1])
		}
	}
}

func (g *GARET) AccGroupRelocation() {
	// calculate predicted amount of gas for each account group
	for i := 0; i < g.AccountGroups; i++ {
		for j := 0; j < g.Ncol; j++ {
			w := float64(2 * (j+1)) / float64(g.Ncol * (g.Ncol + 1))
			g.GasPredAcc[i] += int(float64(g.GasUsedAcc[i][j]) * w)
		}
	}

	Nag := g.AccountGroups

	Qa := utils.NewPriorityQueue()
	for i := 0; i < Nag; i++ {
		Qa.Insert(i, float64(g.GasPredAcc[i]))
	}

	for Qa.Len() != 0 {
		acc, _ := Qa.Pop()
		gasA := g.GasPredAcc[acc.(int)]

		min := g.Context.GasLimit
		minIndex := 0
		for i := 0; i < g.Context.NumberOfShards; i++ {
			if g.GasPredShard[i] < min {
				min = g.GasPredShard[i]
				g.GasPredShard[i] += gasA
				minIndex = i
			}
		}

		g.MappingTable[acc.(int)] = minIndex
	}

	// initialize for next relocation cycle
	for i := 0; i < g.Context.NumberOfShards; i++ {
		g.GasPredShard[i] = 0
	}

	for i := 0; i < g.AccountGroups; i++ {
		g.GasPredAcc[i] = 0
		for j := 0; j < g.Ncol; j++ {
			g.GasUsedAcc[i][j] = 0
		}
	}
}