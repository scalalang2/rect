package utils

import "testing"

func TestGetShardAtSaccAddress(t *testing.T) {
	address1 := "0xC93f2250589a6563f5359051c1eA25746549f0D8"
	address2 := "0x4BD5f0Ee173C81d42765154865EE69361b6aD189"
	address3 := "0x5DF9B87991262F6BA471F09758CDE1c0FC1De734"
	shard1 := GetShardSaccAddress(address1)
	shard2 := GetShardSaccAddress(address2)
	shard3 := GetShardSaccAddress(address3)
}