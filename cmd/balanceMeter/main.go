package main

import "github.com/scalalang2/load-balancing-simulator/balancer"

func main() {
	context := balancer.ExpContext{
		FromBlock: 6500000,
		CollationCycle: 100,
		BlockEpoch: 20,
		NumberOfShards: 20,
		GasLimit: 12000000,
		GasCrossShardTx: 21000,
	}

	balanceMeter := balancer.BalanceMeter {
		Context: context,
		AccountGroups: 20,
		Ncol: 5,
	}

	balanceMeter.Init()
	balanceMeter.StartExperiment()
	balanceMeter.SaveToCSV("balanceMeter.csv")
}
