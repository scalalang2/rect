package main

import (
	"github.com/scalalang2/load-balancing-simulator/balancer"
)

func GetUtilization() {

}

func GetMakespan() {

}

func main() {
	saccBalancer := balancer.SACC{
		FromBlock: 6000000,
		ToBlock: 7000000,
		BlockEpoch: 20,
		NumberOfShards: 20,
		GasLimit: 10000000,
	}

	saccBalancer.Init()
	saccBalancer.StartExperiment()
}
