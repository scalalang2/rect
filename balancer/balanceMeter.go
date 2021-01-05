package balancer

import (
	"context"
	"fmt"
	"github.com/scalalang2/load-balancing-simulator/storage"
	"github.com/scalalang2/load-balancing-simulator/utils"
	"strconv"
)

type BalanceMeter struct {
	Context ExpContext
	AccountGroups int
	Ncol int
	CollationUtils [][]BalanceInfo
	CSTxTable [][]int
	GasUsedAcc [][]int
	GasPredAcc []int
	GasPredShard []int
	MappingTable map[int]int // key: i-th account group. value: j-th shard
}

type CSTx struct {
	GroupA int
	GroupB int
}

func (b *BalanceMeter) Init() {
	b.CollationUtils = make([][]BalanceInfo, b.Context.CollationCycle)
	b.MappingTable = make(map[int]int, b.AccountGroups)

	// init collation utils
	for i := 0; i < b.Context.CollationCycle; i++ {
		b.CollationUtils[i] = make([]BalanceInfo, b.Context.NumberOfShards)
	}

	// init cstx table
	b.CSTxTable = make([][]int, b.AccountGroups)
	for i := 0; i < b.AccountGroups; i++ {
		b.CSTxTable[i] = make([]int, b.AccountGroups)
	}

	// init acc gas used
	b.GasUsedAcc = make([][]int, b.AccountGroups)
	for i := 0; i < b.AccountGroups; i++ {
		b.GasUsedAcc[i] = make([]int, b.Ncol)
	}

	// init gas pred acc
	b.GasPredAcc = make([]int, b.AccountGroups)
	b.GasPredShard = make([]int, b.Context.NumberOfShards)

	for i := 0; i < b.AccountGroups; i++ {
		b.MappingTable[i] = i % b.Context.NumberOfShards
	}
}

func (b *BalanceMeter) StartExperiment() {
	db := storage.OpenDatabase()
	var tx storage.Transaction

	utilNumber := 0
	toBlock := b.Context.FromBlock + (b.Context.CollationCycle * b.Context.BlockEpoch)
	for blockNumber := b.Context.FromBlock; blockNumber < toBlock; blockNumber++ {
		cursor := db.FetchTransactions(int64(blockNumber))
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&tx)
			utils.CheckError(err, "failed to decode transactions")

			toAccGroup := utils.GetAccountGroup(tx.ToAddress, b.AccountGroups)
			fromAccGroup := utils.GetAccountGroup(tx.Sender, b.AccountGroups)
			toShardNum := b.MappingTable[toAccGroup]
			fromShardNum := b.MappingTable[fromAccGroup]

			// util 번 째 블록의 to, from 샤드
			toShard := &b.CollationUtils[utilNumber][toShardNum]
			fromShard := &b.CollationUtils[utilNumber][fromShardNum]

			notLimited := toShard.GasUsed + tx.GasUsed < int64(b.Context.GasLimit)

			// consider cross-shard tx
			if utilNumber > 0 {
				prevShard := b.CollationUtils[utilNumber-1][toShardNum]
				notLimited = toShard.GasUsed + tx.GasUsed + int64(prevShard.CrossShards * b.Context.GasCrossShardTx) < int64(b.Context.GasLimit)
			}

			if notLimited {
				b.GasUsedAcc[toAccGroup][utilNumber % b.Ncol] += int(tx.GasUsed)
				toShard.GasUsed += tx.GasUsed
				toShard.Transactions += 1

				if toShardNum != fromShardNum && tx.GasUsed == 21000 {
					b.CSTxTable[fromAccGroup][toAccGroup] += 1
					toShard.CrossShards += 1
					fromShard.CrossShards += 1
					toShard.Transactions += 1
					fromShard.Transactions += 1
				}
			}
		}

		if (blockNumber + 1) % b.Context.BlockEpoch == 0 {
			utilNumber += 1
			if utilNumber % (b.Ncol+1) == 0 {
				b.PrintUtilization(utilNumber-1)
				b.AccGroupRelocation()
			}
		}
	}
}

func (b *BalanceMeter) PrintCSTxTable() {
	fmt.Println("---------------")
	Q := utils.NewPriorityQueue()

	for i := 0; i < b.AccountGroups; i++ {
		for j := 0; j < b.AccountGroups; j++ {
			var cstx CSTx
			cstx.GroupA = i
			cstx.GroupB = j
			Q.Insert(cstx, float64(b.CSTxTable[i][j]))
		}
	}

	for Q.Len() != 0 {
		cstx, _ := Q.Pop()
		cstxT := cstx.(CSTx)
		fmt.Printf("accA: %d, accB: %d = %d\n", cstxT.GroupA, cstxT.GroupB, b.CSTxTable[cstxT.GroupA][cstxT.GroupB])
	}
}

func (b *BalanceMeter) PrintUtilization(utilNumber int) {
	fmt.Printf("%d %v\n", utilNumber, b.CollationUtils[utilNumber])
}

func (b *BalanceMeter) SaveToCSV(filename string) {
	data := make([][]string, b.Context.CollationCycle + 1)
	for i := 0; i < b.Context.CollationCycle + 1; i++ {
		data[i] = make([]string, b.Context.NumberOfShards)
	}

	for j := 0; j < b.Context.NumberOfShards; j++ {
		data[0][j] = "shard " + strconv.Itoa(j)
	}

	for i := 1; i <= b.Context.CollationCycle; i++ {
		for j := 0; j < b.Context.NumberOfShards; j++ {
			data[i][j] = strconv.FormatInt(b.CollationUtils[i-1][j].GasUsed, 10)
		}
	}

	utils.ReportToCSV(filename, data)
}

func (b *BalanceMeter) AccGroupRelocation() {
	Nag := b.AccountGroups
	visited := make([]int, Nag)
	Qcstx := utils.NewPriorityQueue()

	for i := 0; i < Nag; i++ {
		for j := 0; j < Nag; j++ {
			var cstx CSTx
			cstx.GroupA = i
			cstx.GroupB = j
			priority := (b.CSTxTable[i][j] + b.CSTxTable[j][i]) * b.Context.GasCrossShardTx
			Qcstx.Insert(cstx, float64(priority))
		}
	}

	for Qcstx.Len() != 0 {
		_cstx, _ := Qcstx.Pop()

		var cstx CSTx
		var cstx_gas int
		cstx = _cstx.(CSTx)
		cstx_gas = (b.CSTxTable[cstx.GroupA][cstx.GroupB] + b.CSTxTable[cstx.GroupB][cstx.GroupA]) * b.Context.GasCrossShardTx

		minShard := b.FindMinShard()

		if visited[cstx.GroupA] == 0 && visited[cstx.GroupB] == 0 {
			b.GasPredShard[minShard] += cstx_gas
			b.MappingTable[cstx.GroupA] = minShard
			b.MappingTable[cstx.GroupB] = minShard
			visited[cstx.GroupA] = 1
			visited[cstx.GroupB] = 1
		}
	}

	// initialize for next relocation cycle
	b.ClearForRelocation()
}

func (b *BalanceMeter) ClearForRelocation() {
	for i := 0; i < b.Context.NumberOfShards; i++ {
		b.GasPredShard[i] = 0
	}

	for i := 0; i < b.AccountGroups; i++ {
		b.GasPredAcc[i] = 0
		for j := 0; j < b.Ncol; j++ {
			b.GasUsedAcc[i][j] = 0
		}
	}

	for i := 0; i < Nag; i++ {
		for j := 0; j < Nag; j++ {
			b.CSTxTable[i][j] = 0
			b.CSTxTable[j][i] = 0
		}
	}
}

func (b *BalanceMeter) FindMinShard() int {
	min := 1 << 32 - 1 // max value of integer.
	minIndex := 0
	for i := 0; i < b.Context.NumberOfShards; i++ {
		if b.GasPredShard[i] < min {
			min = b.GasPredShard[i]
			minIndex = i
		}
	}

	return minIndex
}