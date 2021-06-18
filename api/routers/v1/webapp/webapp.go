package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/unknwon/com"
)

// @Summary product detail
// @Tags Product
// @Description 详情
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [get]
func Detail(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface {}
	if db.ExistWebappById(id) {
		data = db.GetWebapp(id)
		c.JSON(msg.SuccessResp(data))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}


// @Summary product list
// @Tags Product
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.WebappSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/ [get]
func Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := db.WebappSearchField{Search: ""}
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}

	apps := db.GetWebapps(page, pageSize, &field)
	data["data"] = apps
	total := db.GetWebappsTotal(&field)
	data["total"] = total
	c.JSON(msg.SuccessResp(data))
	return
}

// @Summary product add
// @Tags Product
// @Description 新增
// @Produce  json
// @Security token
// @Param plugin body db.Webapp true "webapp"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/ [post]
func Add(c *gin.Context) {
	app := db.Webapp{}
	err := c.BindJSON(&app)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	if db.ExistWebappByName(app.Name){
		c.JSON(msg.ErrResp("漏洞名称已存在"))
		return
	} else {
		db.AddWebapp(app)
		c.JSON(msg.SuccessResp(app))
		return
	}
}

// @Summary product update
// @Tags Product
// @Description 更新
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Param webapp body db.Webapp true "webapp"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [put]
func Update(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	app := db.Webapp{}
	err := c.ShouldBindJSON(&app)
	if err != nil {
		c.JSON(msg.ErrResp("组件名称不可为空"))
		return
	}

	if db.ExistWebappById(id){
		db.EditWebapp(id, app)
		c.JSON(msg.SuccessResp(app))
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}

// @Summary product delete
// @Tags Product
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [delete]
func Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	if db.ExistWebappById(id) {
		db.DeleteWebapp(id)
		c.JSON(msg.SuccessResp("删除成功"))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}

}