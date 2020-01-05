package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	Id    string                 `bson:"_id" json:"id"`
	Title string                 `bson:"title" json:"title" binding:"required"`
	Type  string                 `bson:"type" json:"type" binding:"required"`
	Start int64                  `bson:"start" json:"start"`
	End   int64                  `bson:"end" json:"end"`
	Info  map[string]interface{} `bson:"info" json:"info" binding:"required"`
}

var taskCollection *mongo.Collection

func init() {
	taskCollection = model.Collection("task")
}

func (t *Task) Insert() bool {
	_, err := taskCollection.InsertOne(context.TODO(), t)
	return err == nil
}

func (t *Task) Get() bool {
	filter := bson.D{{"_id", t.Id}}
	err := taskCollection.FindOne(context.TODO(), filter).Decode(t)
	return err == nil
}

func (t *Task) GetList(list []string) (results []Task) {
	cur, err := taskCollection.Find(context.TODO(), bson.D{{"_id", bson.M{"$in": list}}}, nil)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		var elem Task
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, elem)
	}
	_ = cur.Close(context.TODO())
	return
}

func (t *Task) Update() bool {
	_, err := taskCollection.UpdateOne(context.TODO(), bson.M{"_id": t.Id}, bson.M{"$set": t})
	return err == nil
}

func (t *Task) Delete() bool {
	_, err := taskCollection.DeleteOne(context.TODO(), bson.M{"_id": t.Id})
	return err == nil
}
