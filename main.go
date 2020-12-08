package main

import (
	"context"
	"fmt"
	"load-balancing-simulator/storage"
)

func main() {
	// averaging formula: avg += (x - avg) / i;
	var opcode storage.Opcode
	var avg float64
	var count int64

	db := storage.OpenDatabase()
	cursor := db.FindOpcodes(0, 1000)
	for cursor.Next(context.TODO()) {
		cursor.Decode(&opcode)

		count++
		avg += (float64(opcode.ElapsedTime) - avg) / float64(count)
	}

	fmt.Println(avg)
}