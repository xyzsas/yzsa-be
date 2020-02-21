package services

import (
	"unicode/utf8"
	"yzsa-be/models"
)

type form struct{}

func init() {
	Index["Form"] = new(form)
}

func (*form) Open(t *models.Task) bool {
	return true
}

func (*form) Realtime(t *models.Task) interface{} {
	return t.Info
}

func (*form) Response(t *models.Task, userId string, resp map[string]interface{}) (code int, reason string) {
	r := &models.Record{Id: t.Id}
	rec := make(map[string]interface{}, 0)
	if len(t.Info) != len(resp) {
		return 400, "表单数据错误"
	}
	for k, v := range t.Info { // k:题号 v:问题对象
		if resp[k] == nil {
			return 400, "表单数据错误"
		}
		question := v.(map[string]interface{}) // 当前问题对象
		switch question["type"].(string) {
		case "choose": // 选择题
			chosen, ok := resp[k].([]interface{}) // 选择的选项编号列表
			if !ok {
				return 400, "表单数据错误"
			}
			if len(chosen) < int(question["min"].(float64)) || // 选项个数不在范围内
				len(chosen) > int(question["max"].(float64)) {
				return 400, "表单数据错误"
			}
			cm := question["choice"].(map[string]interface{})
			for _, ci := range chosen { // 遍历每一个选项
				c, ok := ci.(string)
				// 选项标号不合法
				if !ok || cm[c] == nil {
					return 400, "表单数据错误"
				}
			}
			rec[k] = chosen
		case "fill": // 填空题
			content, ok := resp[k].(string)
			if !ok || utf8.RuneCountInString(content) > int(question["max"].(float64)) {
				return 400, "表单数据错误"
			}
			rec[k] = content
		default:
			return 400, "表单数据错误"
		}
	}
	if r.AddRecord(userId, rec) {
		return 200, "成功"
	} else {
		return 500, "服务器错误"
	}
}
