package middleware

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/services"
	"yzsa-be/utils"
)

func Task(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Keys["id"].(string)
		task := c.Param("id")
		t := &models.Task{Id: task}
		if !t.Get() {
			c.Abort()
			c.String(404, "任务不存在")
			return
		}
		c.Keys["task"] = t
		bypass := utils.Token.Check(c.Keys["id"].(string)+"_"+task, "OVERRIDE")
		if c.Keys["user"] == nil {
			u := &models.User{Id: id}
			if !u.Get() {
				c.Abort()
				c.String(403, "用户信息不存在")
				return
			}
			c.Keys["user"] = u
		}
		if !bypass && !services.Task.CheckPermission(c.Keys["user"].(*models.User).Permission, task) {
			c.Abort()
			c.String(403, "权限不足")
			return
		}
		if !admin && !utils.Cache.Exist(task) {
			c.Abort()
			c.String(409, "任务未开放")
			return
		}
		if admin || bypass {
			c.Next()
			return
		}
		now := utils.TimeStamp(0)
		if now < t.Start || now > t.End {
			c.Abort()
			c.String(403, "任务不在开放时间")
			return
		}
		c.Next()
	}
}
