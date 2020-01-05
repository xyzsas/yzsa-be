package services

import (
	"yzsa-be/models"
	"yzsa-be/utils"
)

var Index map[string]interface {
	Open(*models.Task) bool
	Response(*models.Task, string, map[string]interface{}) bool
}

type task struct{}

func (*task) TaskList(id string, open bool) (results []models.Task) {
	var taskMap map[string]bool
	taskMap = make(map[string]bool)
	u := &models.User{Id: id}
	if !u.Get() {
		return
	}
	permission := u.Permission
	for permission != "" {
		p := &models.Permission{Id: permission}
		if !p.Get() {
			return
		}
		for _, v := range p.Tasks {
			if !taskMap[v] {
				t := &models.Task{Id: v}
				if t.Get() {
					if !open || t.Start != 0 {
						taskMap[v] = true
						results = append(results, *t)
					}
				}
			}
		}
		permission = p.Father
	}
	return
}

func (*task) CheckPermission(permission, task string) bool {
	for permission != "" {
		p := &models.Permission{Id: permission}
		if !p.Get() {
			return false
		}
		for _, v := range p.Tasks {
			if v == task {
				return true
			}
		}
		permission = p.Father
	}
	return false
}

func (*task) Open(t *models.Task) bool {
	if Index[t.Type] == nil {
		return false
	}
	return Index[t.Type].Open(t)
}

func (*task) Response(t *models.Task, userId string, resp map[string]interface{}) (res bool) {
	if Index[t.Type] == nil {
		return false
	}
	res = Index[t.Type].Response(t, userId, resp)
	if res {
		utils.Cache.SAdd(t.Id+"_record", userId)
	}
	return
}

var Task = new(task)
