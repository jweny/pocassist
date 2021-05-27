package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/unknwon/com"
)

//获取 webapp
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

//新增
func Create(c *gin.Context) {
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