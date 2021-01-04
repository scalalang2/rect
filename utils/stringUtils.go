package utils

import (
	"strconv"
)

func GetShardSaccAddress(address string, numShards int) int {
	slices := []rune(address)
	prefix := string(slices[2:4])
	num, _ := strconv.ParseInt(prefix, 16, 32)
	return int(num) % numShards
}