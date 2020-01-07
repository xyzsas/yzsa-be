/* YZSA - Notice
 * Info: content: [text]
 * Record: userid: true
 */

package services

import (
	"fmt"
	"yzsa-be/models"
	"yzsa-be/utils"
)

type notice struct{}

func init() {
	Index["Notice"] = new(notice)
}

func (*notice) Open(t *models.Task) bool {
	fmt.Println(t)
	return utils.Cache.HSetMany(t.Id, t.Info)
}

func (*notice) Response(t *models.Task, userId string, resp map[string]interface{}) (code int, reason string) {
	return 200, "成功"
}
