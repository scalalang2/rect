package balancer

type ExpContext struct {
	FromBlock int
	CollationCycle int
	BlockEpoch int
	NumberOfShards int
	GasLimit int
}

type BalanceInfo struct {
	GasUsed int64
	ElapsedTime int64
	Transactions int
	CrossShards int
}