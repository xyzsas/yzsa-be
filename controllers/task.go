package controllers

import (
	"github.com/gin-gonic/gin"
	"yzsa-be/models"
	"yzsa-be/services"
	"yzsa-be/utils"
)

type task struct{}

func (*task) GetList(c *gin.Context) {
	res := services.Task.TaskList(c.Keys["id"].(string), false)
	for i := 0; i < len(res); i++ {
		res[i].Info = nil
	}
	c.JSON(200, res)
}

func (*task) GetOpen(c *gin.Context) {
	res := services.Task.TaskList(c.Keys["id"].(string), true)
	for i := 0; i < len(res); i++ {
		res[i].Info = nil
	}
	c.JSON(200, res)
}

func (*task) GetOne(c *gin.Context) {
	c.JSON(200, c.Keys["task"].(*models.Task))
}

func (*task) GetRuntime(c *gin.Context) {
	t := c.Keys["task"].(*models.Task)
	res := services.Task.Realtime(t)
	count := utils.Cache.SCard(t.Id + "_record")
	c.JSON(200, gin.H{
		"count":    count,
		"realTime": res,
	})
}

func (*task) GetRealTime(c *gin.Context) {
	res := services.Task.Realtime(c.Keys["task"].(*models.Task))
	c.JSON(200, res)
}

func (*task) Insert(c *gin.Context) {
	var input models.Task
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要title, type, info")
		return
	}
	input.Id = utils.GetRandomString(16)
	if !input.Insert() {
		c.String(500, "服务器错误")
		return
	}
	if (&models.Permission{Id: c.Keys["user"].(*models.User).Permission}).AddTask(input.Id) &&
		(&models.Permission{Id: "admin"}).AddTask(input.Id) &&
		(&models.Record{Id: input.Id, Records: make(map[string]interface{}, 0)}).Insert() {
		c.String(200, "添加成功")
	} else {
		c.String(403, "权限节点不存在")
	}
}

func (*task) Update(c *gin.Context) {
	id := c.Param("id")
	var input models.Task
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要type, title, info")
		return
	}
	if c.Keys["task"].(*models.Task).Start != 0 {
		c.String(409, "任务正在进行中")
		return
	}
	input.Id = id
	if !input.Update() {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, "修改成功")
}

func (*task) Delete(c *gin.Context) {
	id := c.Param("id")
	if c.Keys["task"].(*models.Task).Start != 0 {
		c.String(409, "任务正在进行中")
		return
	}
	t := c.Keys["task"].(*models.Task)
	if t.Delete() &&
		(&models.Permission{}).DeleteTask(id) &&
		(&models.Record{Id: id}).Delete() {
		c.String(200, "删除成功")
	} else {
		c.String(500, "服务器错误")
	}

}

func (*task) Open(c *gin.Context) {
	var input struct {
		Start int64 `json:"start" binding:"required"`
		End   int64 `json:"end" binding:"required"`
	}
	if c.ShouldBind(&input) != nil {
		c.String(400, "参数错误，需要start, end")
		return
	}
	t := c.Keys["task"].(*models.Task)
	if t.Start != 0 {
		c.String(409, "任务正在进行中")
		return
	}
	t.Start = input.Start
	t.End = input.End
	if !t.Update() {
		c.String(500, "服务器错误")
		return
	}
	if !services.Task.Open(c.Keys["task"].(*models.Task)) {
		c.String(500, "服务器错误")
		return
	}
	c.String(200, "任务开启成功")
}

func (*task) Close(c *gin.Context) {
	id := c.Param("id")
	t := c.Keys["task"].(*models.Task)
	t.Start = 0
	t.End = 0
	if t.Update() && utils.Cache.Delete(id) && utils.Cache.Delete(id+"_record") {
		c.String(200, "任务关闭成功")
	} else {
		c.String(500, "任务关闭失败")
	}
}

func (*task) GetUser(c *gin.Context) {
	id := c.Param("id")
	permissions := make([]string, 0)
	for _, v := range (&models.Permission{}).GetByTask(id) {
		permissions = append(permissions, v.Id)
	}
	for i := 0; i < len(permissions); i++ {
		cur := &models.Permission{Id: permissions[i]}
		if !cur.Get() {
			continue
		} else {
			for _, v := range cur.GetChildren() {
				permissions = append(permissions, v.Id)
			}
		}
	}
	res := (&models.User{}).GetGroup(permissions)
	if res == nil {
		c.String(500, "服务器错误")
	} else {
		c.JSON(200, res)
	}
}

var Task = new(task)
