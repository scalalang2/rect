package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"load-balancing-simulator/utils"
	"log"
	"time"
)

type Opcode struct {
	BlockNumber int64 `json:”blockNumber,omitempty”`
	Sender string `json:”sender,omitempty”`
	ContractAddress string `json:”contractAddress,omitempty”`
	Opcode string `json:”opcode,omitempty”`
	ElapsedTime int64 `json:”elapsedTime,omitempty”`
}

func main() {
	dbUser := utils.GetEnv("MONGODB_USER")
	dbPassword := utils.GetEnv("MONGODB_PASSWORD")
	dbURI := utils.GetEnv("MONGODB_URL")
	serverUrl := fmt.Sprintf("mongodb://%s:%s@%s:27017/", dbUser, dbPassword, dbURI)

	client, err := mongo.NewClient(options.Client().ApplyURI(serverUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	var opcode Opcode
	opcodesCol := client.Database("balanceMeter").Collection("opcodes")
	options := options.Find()
	options.SetLimit(10)
	cursor, _ := opcodesCol.Find(context.TODO(), bson.D{}, options)

	for cursor.Next(context.TODO()) {
		cursor.Decode(&opcode)
		fmt.Println(opcode)
	}
}