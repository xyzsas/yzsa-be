package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"yzsa-be/utils"
)

var conn *mongo.Client
var model *mongo.Database

func init() {
	clientOptions := options.Client().ApplyURI(utils.Config.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	conn = client
	model = client.Database(utils.Config.DatabaseName)

	go checkMongoDBAlive()
}

func checkMongoDBAlive() {
	for {
		err := conn.Ping(context.TODO(), nil)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "MongoDB Alive")
		}
		time.Sleep(time.Minute * 30)
	}
}
