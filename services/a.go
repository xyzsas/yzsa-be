package services

import "yzsa-be/models"

func init() {
	Index = make(map[string]interface {
		Open(*models.Task) bool
		Realtime(*models.Task) interface{}
		Response(*models.Task, string, map[string]interface{}) (code int, reason string)
	})
}
