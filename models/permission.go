package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Permission struct {
	Id     string   `bson:"_id" json:"id" binding:"required"`
	Name   string   `bson:"name" json:"name"`
	Father string   `bson:"father" json:"father" binding:"required"`
	Tasks  []string `bson:"tasks" json:"tasks"`
}

var permissionCollection *mongo.Collection

func init() {
	permissionCollection = model.Collection("permission")
	people := &Permission{
		Id:     "people",
		Name:   "用户",
		Father: "",
		Tasks:  make([]string, 0),
	}
	admin := &Permission{
		Id:     "admin",
		Name:   "管理员",
		Father: "",
		Tasks:  make([]string, 0),
	}
	student := &Permission{
		Id:     "student",
		Name:   "学生",
		Father: "people",
		Tasks:  make([]string, 0),
	}
	teacher := &Permission{
		Id:     "teacher",
		Name:   "教师",
		Father: "people",
		Tasks:  make([]string, 0),
	}
	if !people.Get() {
		people.Insert()
	}
	if !admin.Get() {
		admin.Insert()
	}
	if !student.Get() {
		student.Insert()
	}
	if !teacher.Get() {
		teacher.Insert()
	}
}

func (p *Permission) Insert() bool {
	_, err := permissionCollection.InsertOne(context.TODO(), p)
	return err == nil
}

func (p *Permission) Get() bool {
	filter := bson.D{{"_id", p.Id}}
	err := permissionCollection.FindOne(context.TODO(), filter).Decode(p)
	return err == nil
}

func (p *Permission) getMany(filter bson.D) (results []*Permission) {
	cur, err := permissionCollection.Find(context.TODO(), filter, nil)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		var elem Permission
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		results = append(results, &elem)
	}
	_ = cur.Close(context.TODO())
	return
}

func (p *Permission) GetAll() []*Permission {
	return p.getMany(bson.D{{}})
}

func (p *Permission) GetChildren() []*Permission {
	return p.getMany(bson.D{{"father", p.Id}})
}

func (p *Permission) Update() bool {
	_, err := permissionCollection.UpdateOne(context.TODO(), bson.M{"_id": p.Id}, bson.M{"$set": p})
	return err == nil
}

func (p *Permission) AddTask(task string) bool {
	_, err := permissionCollection.UpdateOne(context.TODO(), bson.M{"_id": p.Id}, bson.M{"$addToSet": bson.M{"tasks": task}})
	return err == nil
}

func (p *Permission) Delete() bool {
	_, err := permissionCollection.DeleteOne(context.TODO(), bson.M{"_id": p.Id})
	return err == nil
}

func (p *Permission) DeleteList(list []string) bool {
	_, err := permissionCollection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": list}}, nil)
	return err == nil
}

func (p *Permission) DeleteTask(task string) bool {
	_, err := permissionCollection.UpdateMany(context.TODO(), bson.M{"tasks": bson.M{"$regex": task}}, bson.M{"$pull": bson.M{"tasks": task}})
	return err == nil
}

func (p *Permission) GetByTask(task string) []*Permission {
	return p.getMany(bson.D{{"tasks", bson.M{"$regex": task}}})
}