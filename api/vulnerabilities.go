package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pocassist/database"
)

type VulSerializer struct {
	// 返回给前端的字段
	WriterName		string						`json:"writer_name"`
	WebAppName		string						`json:"webapp_name"`
	Id       	int         `gorm:"primary_key" json:"id"`
	NameZh 		string  `gorm:"column:name_zh" json:"name_zh"`
	Cve 		string     `gorm:"column:cve" json:"cve"`
	Cnnvd 		string   `gorm:"column:cnnvd" json:"cnnvd"`
	Severity	string    `gorm:"column:severity" json:"severity"`
	Category 	string    `gorm:"column:category" json:"category"`
	Description string    `gorm:"column:description" json:"description"`
	Suggestion	string  `gorm:"column:suggestion" json:"suggestion"`
	Language 	string    `gorm:"column:language" json:"language"`
	Webapp		int      `gorm:"column:webapp" json:"webapp"`
}

//获取 webapp
func GetWebApps(c *gin.Context) {
	data := make(map[string]interface{})
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	apps := database.GetWebApps(page, pageSize)
	data["data"] = apps
	total := database.GetWebAppsTotal()
	data["total"] = total
	c.JSON(SuccessResp(data))
	return
}

//新增
func CreateWebApp(c *gin.Context) {
	app := database.Webapp{}
	err := c.BindJSON(&app)
	if err != nil {
		c.JSON(ErrResp("参数校验不通过"))
		return
	}
	if database.ExistWebappByName(app.Name){
		c.JSON(ErrResp("漏洞名称已存在"))
		return
	} else {
		database.AddWebapp(app)
		c.JSON(SuccessResp(app))
		return
	}
}

//获取单个描述
func GetVul(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface {}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if ! valid.HasErrors() {
		if database.ExistVulnerabilityByID(id) {
			data = database.GetVulnerability(id)
			c.JSON(SuccessResp(data))
			return
		} else {
			c.JSON(ErrResp("record not found"))
			return
		}
	} else {
		c.JSON(ErrResp(DealValidError(valid)))
		return
	}
}

//获取多个描述
func GetVuls(c *gin.Context) {
	data := make(map[string]interface{})
	field := database.VulnerabilitySearchField{
		Search:"",
		WebappField:-1,
		CategoryField:"",}

	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}
	if arg := c.Query("productField"); arg != "" {
		appId := com.StrTo(arg).MustInt()
		field.WebappField = appId
	}
	if arg := c.Query("typeField"); arg != "" {
		field.CategoryField = arg
	}

	vuls := database.GetVulnerabilities(page, pageSize, &field)
	var vulRespData []VulSerializer
	// 获取上传者
	token := c.Request.Header.Get("Authorization")
	claims, _ := ParseToken(token)

	for _, vul := range vuls {
		var appName string
		if vul.ForeignWebapp != nil {
			appName = vul.ForeignWebapp.Name
		} else {
			appName = ""
		}
		vulRespData = append(vulRespData, VulSerializer{
			WriterName:    claims.Username,
			WebAppName:    appName,
			Id:vul.Id,
			NameZh:vul.NameZh,
			Cve:vul.Cve,
			Cnnvd:vul.Cnnvd,
			Severity:vul.Severity,
			Category:vul.Category,
			Description:vul.Description,
			Suggestion:vul.Suggestion,
			Language:vul.Language,
			Webapp:vul.Webapp,
		})
	}
	data["data"] = vulRespData
	total := database.GetVulnerabilitiesTotal(&field)
	data["total"] = total
	c.JSON(SuccessResp(data))
	return
}

//新增
func CreateVul(c *gin.Context) {
	vul := database.Vulnerability{}
	err := c.BindJSON(&vul)
	if err != nil {
		c.JSON(ErrResp("参数校验不通过"))
		return
	}
	if database.ExistVulnerabilityByNameZh(vul.NameZh){
		c.JSON(ErrResp("漏洞名称已存在"))
		return
	} else {
		database.AddVulnerability(vul)
		c.JSON(SuccessResp(vul))
		return
	}
}

//修改
func UpdateVul(c *gin.Context) {
	vul := database.Vulnerability{}
	err := c.BindJSON(&vul)
	if err != nil {
		c.JSON(ErrResp("参数校验不通过"))
		return
	}

	valid := validation.Validation{}

	valid.Min(vul.Id, 1, "id").Message("ID必须大于0")
	valid.Required(vul.NameZh, "Affects").Message("Affects不能为空")

	if ! valid.HasErrors() {
		if database.ExistVulnerabilityByID(vul.Id){
			database.EditVulnerability(vul.Id, vul)
			c.JSON(SuccessResp(vul))
		} else {
			c.JSON(ErrResp("record not found"))
			return
		}
	} else {
		c.JSON(ErrResp(DealValidError(valid)))
		return
	}
}

//删除
func DeleteVul(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if ! valid.HasErrors() {
		if database.ExistVulnerabilityByID(id) {
			database.DeleteVulnerability(id)
			c.JSON(SuccessResp("删除成功"))
			return
		} else {
			c.JSON(ErrResp("record not found"))
			return
		}
	} else {
		c.JSON(ErrResp(DealValidError(valid)))
		return
	}
}

type Basic struct {
	Name 	string 	`json:"name"`
	Label 	string	`json:"label"`
}

func GetBasic(c *gin.Context) {
	var LANGUAGE_CHOICE []Basic
	for _, v := range []string{"Any","ASP","JAVA","Python","NodeJS","PHP","Ruby","ASPX"} {
		LANGUAGE_CHOICE = append(LANGUAGE_CHOICE, Basic{Name:v, Label:v})
	}
	var AFFECT_CHOICE []Basic
	for _, v := range []string{"server","text","directory","url","appendparam","replaceparam","script"} {
		AFFECT_CHOICE = append(AFFECT_CHOICE, Basic{Name:v, Label:v})
	}
	var LEVEL_CHOICE []Basic
	for _, v := range []string{"high","middle","low","info",} {
		LEVEL_CHOICE = append(LEVEL_CHOICE, Basic{Name:v, Label:v})
	}
	var TYPE_CHOICE []Basic
	for _, v := range []string{"SQL 注入","命令执行","信息泄漏","其他类型","发现备份文件","未知","目录穿越","未授权","ShellCode","任意文件下载","任意文件读取","反序列化","任意文件写入","弱口令","权限提升","目录遍历","JAVA反序列化","代码执行","嵌入恶意代码","拒绝服务","文件上传","远程文件包含","跨站请求伪造","跨站脚本XSS","XPath注入","缓冲区溢出","XML注入","服务器端请求伪造","Cookie验证错误","解析错误","本地文件包含","配置错误"} {
		TYPE_CHOICE = append(TYPE_CHOICE, Basic{Name:v, Label:v})
	}

	data := make(map[string]interface{})
	data["VulLanguage"] = LANGUAGE_CHOICE
	data["VulLevel"] = LEVEL_CHOICE
	data["ModuleAffects"] = AFFECT_CHOICE
	data["VulType"] = TYPE_CHOICE

	c.JSON(SuccessResp(data))
	return
}