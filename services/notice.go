/* YZSA - Notice
 * Info: content: [text]
 * Record: userid: true
 */

package services

import (
	"yzsa-be/models"
)

type notice struct{}

func init() {
	Index["Notice"] = new(notice)
}

func (*notice) Open(t *models.Task) bool {
	return true
}

func (*notice) Realtime(t *models.Task) interface{} {
  return t.Info
}

func (*notice) Response(t *models.Task, userId string, resp map[string]interface{}) (code int, reason string) {
	return 200, "成功"
}
