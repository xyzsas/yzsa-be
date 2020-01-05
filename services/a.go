package services

import "yzsa-be/models"

func init() {
	Index = make(map[string]interface {
		Open(*models.Task) bool
		Response(*models.Task, string, map[string]interface{}) bool
	})
}
