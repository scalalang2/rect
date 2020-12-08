package storage

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

type DB struct {
	client *mongo.Client
}

// open database connection
func OpenDatabase() *DB {
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
	utils.CheckError(err, "Database Client Connection Error")

	var db DB
	db.client = client

	return &db
}

func (db *DB) FindOpcodes(skip int64, limit int64) *mongo.Cursor {
	col := db.client.Database("balanceMeter").Collection("opcodes")
	options := options.Find()
	options.SetSkip(skip)
	options.SetLimit(limit)

	cursor, err := col.Find(context.TODO(), bson.D{}, options)
	utils.CheckError(err, "db.FindOpcodes() Error")

	return cursor
}