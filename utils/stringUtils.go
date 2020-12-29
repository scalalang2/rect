package utils

import (
	"strconv"
)

func GetShardSaccAddress(address string, numShards int) uint64 {
	slices := []rune(address)
	prefix := string(slices[2:4])
	num, _ := strconv.ParseUint(prefix, 16, 32)
	return num % uint64(numShards)
}