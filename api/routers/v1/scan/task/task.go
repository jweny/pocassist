package task

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/unknwon/com"
)

// @Summary task list
// @Tags Task
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.TaskSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/task/ [get]
func Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := db.TaskSearchField{Search: ""}

	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}

	tasks := db.GetTask(page, pageSize, &field)
	data["data"] = tasks

	total := db.GetTaskTotal(&field)
	data["total"] = total

	c.JSON(msg.SuccessResp(data))
	return
}


// @Summary task delete
// @Tags Task
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/task/{id}/ [delete]
func Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	if db.ExistTaskByID(id) {
		db.DeleteTask(id)
		c.JSON(msg.SuccessResp("删除成功"))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}

}