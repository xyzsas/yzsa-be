package controllers

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/utils"
)

type user struct{}

func (*user) HandShake(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.String(400, "参数错误，需要id")
		return
	}
	token := utils.Token.SetHandShake(id)
	if token == "" {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, token)
}

func (*user) Login(c *gin.Context) {
	var input struct {
		Id       string `json:"id" binding:"required"`
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要id, token, password")
		return
	}
	token, expire := utils.Token.GetHandShake(input.Id)
	if token != input.Token {
		c.String(403, "上下文错误")
		return
	}
	if utils.TimeStamp(0) <= expire {
		c.String(425, "您输入密码的速度太快啦！")
		return
	}
	u := &models.User{Id: input.Id}
	if !u.Get() || utils.HASH(u.Password, token) != input.Password {
		c.String(403, "账户名或密码错误")
		return
	}
	token = utils.Token.Set(u.Id)
	if token == "" {
		c.String(500, "服务器错误")
		return
	}
	c.Header("token", token)
	c.JSON(200, u)
}

func (*user) Logout(c *gin.Context) {
	id := c.GetHeader("id")
	c.Header("token", "")
	if utils.Token.Delete(id) {
		c.String(200, "登出成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*user) UpdatePwd(c *gin.Context) {
	var input struct {
		OldPwd string `json:"oldPwd" binding:"required"`
		NewPwd string `json:"newPwd" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要oldPwd, newPwd")
		return
	}
	u := &models.User{Id: c.Keys["id"].(string)}
	if !u.Get() {
		c.String(404, "用户不存在")
		return
	}
	if u.Password != utils.HASH(input.OldPwd, utils.Config.Salt) {
		c.String(403, "原密码错误，请重试")
		return
	}
	u.Password = utils.HASH(input.NewPwd, utils.Config.Salt)
	if !u.Update() {
		c.String(500, "服务器错误")
		return
	}
	if utils.Token.Delete(u.Id) {
		c.String(200, "密码修改成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*user) GetList(c *gin.Context) {
	role := c.DefaultQuery("role", "")
	permission := c.DefaultQuery("permission", "")
	if role != "" {
		c.JSON(200, (&models.User{}).GetByRole(role))
		return
	}
	if permission != "" {
		c.JSON(200, (&models.User{}).GetByPermission(permission))
		return
	}
	c.String(400, "参数错误，需要role或permission")
}

func (*user) GetOne(c *gin.Context) {
	id := c.Param("id")
	u := &models.User{Id: id}
	if !u.Get() {
		c.String(404, "该用户不存在")
		return
	}
	c.JSON(200, u)
}

func (*user) Insert(c *gin.Context) {
	var input models.User
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要id, name, role, permission")
		return
	}
	if input.Role == "admin" || input.Permission == "admin" {
		c.String(403, "权限不足")
		return
	}
	if !(&models.Permission{Id: input.Permission}).Get() {
		c.String(403, "权限节点不存在")
		return
	}
	u := &models.User{Id: input.Id}
	if u.Get() {
		c.String(400, "该用户已存在")
		return
	}
	if !input.Insert() {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, "添加成功")
}

func (*user) UpdatePermission(c *gin.Context) {
	var input struct {
		Permission string `json:"permission" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要permission")
		return
	}
	id := c.Param("id")
	u := &models.User{Id: id}
	if !u.Get() {
		c.String(404, "该用户不存在")
		return
	}
	if u.Role == "admin" || input.Permission == "admin" {
		c.String(403, "权限不足")
		return
	}
	u.Permission = input.Permission
	if !u.Update() {
		c.String(500, "服务器错误")
		return
	}
	if utils.Token.Delete(id) {
		c.String(200, "用户权限节点修改成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*user) UpdateName(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要name")
		return
	}
	id := c.Param("id")
	u := &models.User{Id: id}
	if !u.Get() {
		c.String(404, "该用户不存在")
		return
	}
	u.Name = input.Name
	if !u.Update() {
		c.String(500, "服务器错误")
		return
	}
	if utils.Token.Delete(id) {
		c.String(200, "用户名称修改成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*user) Delete(c *gin.Context) {
	id := c.Param("id")
	u := &models.User{Id: id}
	if !u.Get() {
		c.String(404, "该用户不存在")
		return
	}
	if u.Role == "admin" {
		c.String(403, "权限不足")
		return
	}
	if !u.Delete() {
		c.String(500, "服务器错误")
		return
	}
	if utils.Token.Delete(id) {
		c.String(200, "用户删除成功")
	} else {
		c.String(500, "服务器错误")
	}
}

func (*user) ResetPwd(c *gin.Context) {
	id := c.Param("id")
	u := &models.User{Id: id}
	if !u.Get() {
		c.String(404, "该用户不存在")
		return
	}
	if u.Role == "admin" {
		c.String(403, "权限不足")
		return
	}
	u.Password = utils.HASH(id, utils.Config.Salt)
	if !u.Update() {
		c.String(500, "服务器错误")
		return
	}
	if utils.Token.Delete(id) {
		c.String(200, "密码重置成功")
	} else {
		c.String(500, "服务器错误")
	}
}

var User = new(user)
