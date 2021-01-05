package utils

import (
	"fmt"
	"testing"
)

// S-ACC, prefix string을 보고 샤드 넘버를 할당하기 위한 유틸 함수 테스트
func TestGetShardAtSaccAddress(t *testing.T) {
	address1 := "0xC93f2250589a6563f5359051c1eA25746549f0D8"
	address2 := "0x4BD5f0Ee173C81d42765154865EE69361b6aD189"
	address3 := "0x5DF9B87991262F6BA471F09758CDE1c0FC1De734"
	shard1 := GetShardSaccAddress(address1, 20)
	shard2 := GetShardSaccAddress(address2, 20)
	shard3 := GetShardSaccAddress(address3, 20)

	if shard1 != 1 {
		t.Fatalf("failed to parse shard number: %d", shard1)
	}

	if shard2 != 15 {
		t.Fatalf("failed to parse shard number: %d", shard2)
	}

	if shard3 != 13 {
		t.Fatalf("failed to parse shard number: %d", shard3)
	}
}

func TestGetAccountGroup(t *testing.T) {
	address1 := "0xC93f2250589a6563f5359051c1eA25746549f0D8"
	address2 := "0x4BD5f0Ee173C81d42765154865EE69361b6aD189"
	address3 := "0x5DF9B87991262F6BA471F09758CDE1c0FC1De734"
	shard1 := GetAccountGroup(address1, 20)
	shard2 := GetAccountGroup(address2, 20)
	shard3 := GetAccountGroup(address3, 20)
	fmt.Printf("%d %d %d", shard1, shard2, shard3)
}