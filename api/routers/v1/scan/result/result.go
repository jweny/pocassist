package result

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/unknwon/com"
)

// @Summary result list
// @Tags Result
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.ResultSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/result/ [get]
func Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := db.ResultSearchField{Search: "", TaskField: -1, VulField:-1}
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}
	if arg := c.Query("taskField"); arg != "" {
		enable := com.StrTo(arg).MustInt()
		field.TaskField = enable
	}
	if arg := c.Query("vulField"); arg != "" {
		vul := com.StrTo(arg).MustInt()
		field.VulField = vul
	}
	results := db.GetResult(page, pageSize, &field)
	data["data"] = results

	total := db.GetResultTotal(&field)
	data["total"] = total

	c.JSON(msg.SuccessResp(data))
	return
}

// @Summary result delete
// @Tags Result
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/result/{id}/ [delete]
func Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	if db.ExistResultByID(id) {
		db.DeleteResult(id)
		c.JSON(msg.SuccessResp("删除成功"))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}