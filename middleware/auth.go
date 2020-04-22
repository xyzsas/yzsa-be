package middleware

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/utils"
)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("id")
		token := c.GetHeader("token")
		if utils.Token.Check(id, token) {
			c.Keys = make(map[string]interface{})
			c.Keys["id"] = id
			newToken := utils.Token.Set(id)
			if newToken == "" {
				c.Abort()
				c.String(500, "服务器错误，请联系管理员")
				return
			} else {
				c.Header("token", newToken)
			}
			c.Next()
		} else {
			c.Abort()
			c.String(401, "很抱歉，需要验证您的身份")
		}
	}
}

func RoleAuth(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := &models.User{Id: c.Keys["id"].(string)}
		if !u.Get() {
			c.Abort()
			c.String(403, "用户信息不存在")
			return
		}
		if u.Role != role && u.Role != "admin" {
			c.Abort()
			c.String(403, "权限不足")
			return
		}
		c.Keys["user"] = u
		c.Next()
	}
}
