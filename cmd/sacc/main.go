package main

import (
	"github.com/scalalang2/load-balancing-simulator/balancer"
)

func main() {
	context := balancer.ExpContext{
		FromBlock: 6000000,
		CollationCycle: 100,
		BlockEpoch: 20,
		NumberOfShards: 20,
		GasLimit: 10000000,
	}

	saccBalancer := balancer.SACC { Context: context }

	saccBalancer.Init(false)
	saccBalancer.StartExperiment()
	saccBalancer.SaveToCSV("sacc_without_cstx.csv")

	saccBalancer.Init(true)
	saccBalancer.StartExperiment()
	saccBalancer.SaveToCSV("sacc_with_cstx.csv")
}
