package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"yzsa-be/utils"
)

type User struct {
	Id         string `bson:"_id" json:"id" binding:"required"`
	Name       string `bson:"name" json:"name" binding:"required"`
	Role       string `bson:"role" json:"role" binding:"required"`
	Password   string `bson:"password" json:"-"`
	Permission string `bson:"permission" json:"permission" binding:"required"`
}

var userCollection *mongo.Collection

func init() {
	userCollection = model.Collection("user")
	admin := &User{
		Id:         "admin",
		Name:       "超级管理员",
		Role:       "admin",
		Password:   utils.HASH("admin", utils.Config.Salt),
		Permission: "admin",
	}
	if !admin.Get() {
		admin.Insert()
	}
}

func (u *User) Insert() bool {
	_, err := userCollection.InsertOne(context.TODO(), u)
	return err == nil
}

func (u *User) Get() bool {
	filter := bson.D{{"_id", u.Id}}
	err := userCollection.FindOne(context.TODO(), filter).Decode(u)
	return err == nil
}

func (u *User) GetByRole(role string) (results []User) {
	cur, err := userCollection.Find(context.TODO(), bson.D{{"role", role}}, nil)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, elem)
	}
	_ = cur.Close(context.TODO())
	return
}

func (u *User) GetByPermission(permission string) (results []User) {
	cur, err := userCollection.Find(context.TODO(), bson.D{{"permission", permission}}, nil)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, elem)
	}
	_ = cur.Close(context.TODO())
	return
}

func (u *User) Update() bool {
	_, err := userCollection.UpdateOne(context.TODO(), bson.M{"_id": u.Id}, bson.M{"$set": u})
	return err == nil
}

func (u *User) Delete() bool {
	_, err := userCollection.DeleteOne(context.TODO(), bson.M{"_id": u.Id})
	return err == nil
}

func (u *User) DeleteGroup(permissions []string) bool {
	_, err := userCollection.DeleteMany(context.TODO(), bson.M{"permission": bson.M{"$in": permissions}}, nil)
	return err == nil
}
