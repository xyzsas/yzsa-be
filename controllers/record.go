package controllers

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/services"
	"yzsa-be/utils"
)

type record struct{}

func (*record) GetAll(c *gin.Context) {
	id := c.Param("id")
	r := &models.Record{Id: id}
	if !r.Get() {
		c.String(404, "记录不存在")
		return
	}
	t := c.Keys["task"].(*models.Task)
	if t.Start == 0 {
		c.String(403, "任务已关闭")
		return
	}
	res := make(map[string]interface{})
	for k, v := range r.Records {
		if utils.Cache.SExist(id+"_record", k) {
			res[k] = v
		}
	}
	c.JSON(200, res)
}

func (*record) GetOne(c *gin.Context) {
	id := c.Param("id")
	user := c.Keys["id"].(string)
	if !utils.Cache.SExist(id+"_record", user) {
		c.String(404, "记录不存在")
		return
	}
	r := &models.Record{Id: id}
	if !r.GetRecord(user) || len(r.Records) == 0 {
		c.String(200, "任务已完成")
		return
	}
	c.JSON(200, r.Records)
}

func (*record) Response(c *gin.Context) {
	var input struct {
		Data map[string]interface{} `json:"data"`
	}
	id := c.Param("id")
	user := c.Keys["id"].(string)
	if utils.Cache.SExist(id+"_record", user) {
		c.String(403, "不能重复提交")
		return
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要data")
		return
	}
	code, reason := services.Task.Response(c.Keys["task"].(*models.Task), user, input.Data)
	c.String(code, reason)
}

var Record = new(record)
