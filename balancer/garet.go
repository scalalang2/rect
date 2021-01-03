package balancer

type GARET struct {
	Context ExpContext
	AccountGroups int
	Ncol int
	WithCSTx bool
	CollationUtils [][]BalanceInfo
	GasUsedAcc []int
}

func (g GARET) Init(withCSTx bool) {
	g.WithCSTx = withCSTx
	g.CollationUtils = make([][]BalanceInfo, g.Context.CollationCycle)

	// amount of gas used for each account group.
	g.GasUsedAcc = make([]int, g.AccountGroups)

	for i := 0; i < g.Context.CollationCycle; i++ {
		g.CollationUtils[i] = make([]BalanceInfo, g.Context.NumberOfShards)
	}
}

func (g GARET) StartExperiment() {
	// read each transactions

	// calculate account group number

	// get shard num
}