package main

import (
	"context"
	"load-balancing-simulator/storage"
)

func main() {
	// averaging formula: avg += (x - avg) / i;
	var opcode storage.Opcode
	var count int64
	opcodeMap := make(map[string]float64, 50)

	db := storage.OpenDatabase()
	cursor := db.FindOpcodes(0, 1000)
	for cursor.Next(context.TODO()) {
		cursor.Decode(&opcode)
		_, doesExist := opcodeMap[opcode.Opcode]
		if !doesExist {
			opcodeMap[opcode.Opcode] = 0.0
		}

		count++
		opcodeMap[opcode.Opcode] += (float64(opcode.ElapsedTime) - opcodeMap[opcode.Opcode]) / float64(count)
	}
}