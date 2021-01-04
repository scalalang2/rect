package main

import "github.com/scalalang2/load-balancing-simulator/balancer"

func main() {
	context := balancer.ExpContext{
		FromBlock: 6000000,
		CollationCycle: 100,
		BlockEpoch: 20,
		NumberOfShards: 20,
		GasLimit: 12000000,
		GasCrossShardTx: 42000,
	}

	garetBalancer := balancer.GARET {
		Context: context,
		AccountGroups: 100,
		Ncol: 5,
	}

	garetBalancer.Init(false)
	garetBalancer.StartExperiment()
}
