package controllers

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/utils"
)

type permission struct{}

func (*permission) GetList(c *gin.Context) {
	permissions := (&models.Permission{}).GetAll()
	c.JSON(200, permissions)
}

func (*permission) Insert(c *gin.Context) {
	var input models.Permission
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要id, father")
		return
	}
	if input.Father == "admin" {
		c.String(403, "权限不足")
		return
	}
	p := &models.Permission{Id: input.Id}
	if p.Get() {
		c.String(400, "该权限节点已存在")
		return
	}
	fa := &models.Permission{Id: input.Father}
	if !fa.Get() {
		c.String(400, "父节点不存在")
		return
	}
	if len(input.Tasks) == 0 {
		input.Tasks = make([]string, 0)
	}
	if !input.Insert() {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, "添加成功")
}

func (*permission) Update(c *gin.Context) {
	var input models.Permission
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要id, father")
		return
	}
	old := &models.Permission{Id: input.Id}
	if !old.Get() {
		c.String(404, "该权限节点不存在")
		return
	}
	if input.Id == "people" || input.Id == "admin" {
		if input.Father != "" {
			c.String(403, "权限不足")
			return
		}
	} else if input.Id == "student" || input.Id == "teacher" {
		if input.Father != "people" {
			c.String(403, "权限不足")
			return
		}
	} else {
		father := &models.Permission{Id: input.Father}
		if !father.Get() {
			c.String(403, "父节点不存在")
			return
		}
	}
	if len(input.Tasks) == 0 {
		input.Tasks = make([]string, 0)
	}
	if !input.Update() {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, "修改成功")
}

func (*permission) Delete(c *gin.Context) {
	id := c.Param("id")
	p := &models.Permission{Id: id}
	if !p.Get() {
		c.String(404, "该权限节点不存在")
		return
	}
	if id == "people" || id == "admin" || id == "student" || id == "teacher" {
		c.String(403, "权限不足")
		return
	}
	subtree := []string{id}
	for i := 0; i < len(subtree); i++ {
		cur := &models.Permission{Id: subtree[i]}
		if !cur.Get() {
			continue
		} else {
			ch := cur.GetChildren()
			for _, v := range ch {
				subtree = append(subtree, v.Id)
			}
		}
	}
	if !p.DeleteList(subtree) {
		c.String(500, "服务器错误")
		return
	}
	if !(&models.User{}).DeleteGroup(subtree) {
		c.String(500, "删除用户错误")
		return
	}
	c.String(200, "删除权限节点、子节点及所属用户成功")
}

func (*permission) AddTemp(c *gin.Context) {
	taskId := c.Param("id")
	var input struct {
		User       string `json:"user" binding:"required"`
		Expiration int64  `json:"expiration" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要user, expiration")
		return
	}
	if input.Expiration > 86400 {
		c.String(400, "过期时间过长")
		return
	}
	if utils.Token.SetOverride(input.User, taskId, input.Expiration) {
		c.String(200, "成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*permission) DeleteTemp(c *gin.Context) {
	taskId := c.Param("id")
	userId := c.Param("user")
	if utils.Token.Delete(userId + "_" + taskId) {
		c.String(200, "成功")
	} else {
		c.String(500, "服务器错误")
	}
}

var Permission = new(permission)
