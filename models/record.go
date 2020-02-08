package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	Id      string                 `bson:"_id" json:"id" binding:"required"`
	Records map[string]interface{} `bson:"records" json:"records"`
}

var recordCollection *mongo.Collection

func init() {
	recordCollection = model.Collection("record")
}

func (r *Record) Insert() bool {
	_, err := recordCollection.InsertOne(context.TODO(), r)
	return err == nil
}

// 注：请只在需要获取本任务全部记录时使用
func (r *Record) Get() bool {
	filter := bson.M{"_id": r.Id}
	err := recordCollection.FindOne(context.TODO(), filter).Decode(r)
	return err == nil
}

func (r *Record) Delete() bool {
	_, err := recordCollection.DeleteOne(context.TODO(), bson.M{"_id": r.Id})
	return err == nil
}

func (r *Record) AddRecord(id string, data interface{}) bool {
	_, err := recordCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": r.Id},
		bson.M{"$set": bson.M{"records." + id: data}},
	)
	return err == nil
}

func (r *Record) GetRecord(id string) bool {
	opt := options.FindOne()
	opt.Projection = bson.D{{"records." + id, 1}}
	err := recordCollection.FindOne(context.TODO(), bson.D{{"_id", r.Id}}, opt).Decode(r)
	return err == nil && len(r.Records) != 0
}
