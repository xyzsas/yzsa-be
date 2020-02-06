package services

import (
	"go.mongodb.org/mongo-driver/bson"
	"yzsa-be/models"
	"yzsa-be/utils"
)

type courseSelect struct{}

func init() {
	Index["CourseSelect"] = new(courseSelect)
}

func (*courseSelect) Open(t *models.Task) bool {
	return utils.Cache.HSetMany(t.Id, t.Info)
}

func (*courseSelect) Realtime(t *models.Task) interface{} {
	return utils.Cache.HGetAll(t.Id)
}

func (*courseSelect) Response(t *models.Task, userId string, resp map[string]interface{}) (code int, reason string) {
	u := &models.User{Id: userId}
	if !u.Get() || u.Role != "student" {
		return 403, "仅学生可以选课"
	}
	if _, ok := resp["course"]; !ok {
		return 400, "参数错误，需要data.course"
	}
	if _, ok := resp["course"].(string); !ok {
		return 400, "参数错误，需要data.course"
	}
	res := utils.Cache.Run(
		`
			local left = redis.call("HGET", KEYS[1], KEYS[2])
			if(left)
			then
				if(left == "0")
				then
					return 0
				else
					redis.call("HINCRBY", KEYS[1], KEYS[2], -1)
					return 1
				end
			else
				return 0
			end
		`,
		[]string{t.Id, resp["course"].(string)},
	)
	if res.(int64) == 0 {
		return 403, "课程无名额"
	} else {
		r := &models.Record{Id: t.Id}
		if r.AddRecord(userId, bson.M{"course" : resp["course"].(string)}) {
			return 200, "成功"
		} else {
			return 500, "服务器错误"
		}
	}
}
