package services

import (
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

func (*courseSelect) Response(t *models.Task, userId string, resp map[string]interface{}) bool {
	u := &models.User{Id: userId}
	if !u.Get() || u.Role != "student" {
		return false
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
		return false
	} else {
		r := &models.Record{Id: t.Id}
		return r.AddRecord(userId, resp["course"].(string))
	}
}
