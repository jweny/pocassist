package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/unknwon/com"
)

// @Summary product list
// @Tags Product
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/ [get]
func Get(c *gin.Context) {
	data := make(map[string]interface{})
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	apps := db.GetWebApps(page, pageSize)
	data["data"] = apps
	total := db.GetWebAppsTotal()
	data["total"] = total
	c.JSON(msg.SuccessResp(data))
	return
}

// @Summary product add
// @Tags Product
// @Description 新增
// @Produce  json
// @Security token
// @Param plugin body rule.Plugin true "plugin"
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