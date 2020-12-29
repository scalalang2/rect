package main

import (
	"github.com/scalalang2/load-balancing-simulator/balancer"
)

func main() {
	saccBalancer := balancer.SACC{
		FromBlock: 6000000,
		ToBlock: 6001000,
		BlockEpoch: 20,
		NumberOfShards: 20,
		GasLimit: 10000000,
	}

	saccBalancer.Init()
	saccBalancer.StartExperiment()
}
