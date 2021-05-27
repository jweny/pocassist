package v1

import (
	"bufio"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/unknwon/com"
	"gorm.io/datatypes"
	"log"
)

const (
	TargetUrl = "url"
	TargetUrlFile = "file"
	TargetUrlRaw = "raw"
)

type PluginSerializer struct {
	// 返回给前端的字段
	DespName		string		   `json:"desp_name"`
	Id       		int            `gorm:"primary_key" json:"id"`
	VulId           string         `gorm:"column:vul_id" json:"vul_id"`
	Affects         string         `gorm:"column:affects" json:"affects"`
	JsonPoc         datatypes.JSON `gorm:"column:json_poc" json:"json_poc"`
	Enable          bool          `gorm:"column:enable" json:"enable"`
	Description   	int            `gorm:"column:description" json:"description"`
}

type RunSinglePluginSerializer struct {
	// 运行单个
	Target			string		   `json:"target"`
	Affects         string         `gorm:"column:affects" json:"affects"`
	JsonPoc         datatypes.JSON `gorm:"column:json_poc" json:"json_poc"`
}

type RunPluginsSerializer struct {
	// 批量运行
	Target			string		`json:"target"`
	TargetType		string		`json:"target_type"`
	RunType 		string		`json:"run_type"`
	VulIdList		[]string	`json:"vul_id_list"`
}

//获取单个plugin
func GetPlugin(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface {}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if ! valid.HasErrors() {
		if db.ExistPluginByID(id) {
			data = db.GetPlugin(id)
			c.JSON(msg.SuccessResp(data))
			return
		} else {
			c.JSON(msg.ErrResp("record not found"))
			return
		}
	} else {
		c.JSON(msg.ErrResp(msg.DealValidError(valid)))
		return
	}
}

//获取多个pluign
func GetPlugins(c *gin.Context) {
	data := make(map[string]interface{})
	field := db.PluginSearchField{Search: "", EnableField:-1, AffectsField:"",}
	valid := validation.Validation{}
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}
	if arg := c.Query("enableField"); arg != "" {
		enable := com.StrTo(arg).MustInt()
		field.EnableField = enable
		valid.Range(enable, 0, 1, "state").Message("状态只允许0或1")
	}
	if arg := c.Query("affectsField"); arg != "" {
		field.AffectsField = arg
	}
	if ! valid.HasErrors() {
		plugins := db.GetPlugins(page, pageSize, &field)

		var pluginRespData []PluginSerializer

		for _, plugin := range plugins {
			var despName string
			if plugin.Vulnerability != nil {
				despName = plugin.Vulnerability.NameZh
			} else {
				despName = ""
			}

			pluginRespData = append(pluginRespData, PluginSerializer{
				DespName: despName,
				Id: plugin.Id,
				VulId: plugin.VulId,
				Affects: plugin.Affects,
				JsonPoc: plugin.JsonPoc,
				Enable: plugin.Enable,
				Description: plugin.Desc,
			})
		}
		data["data"] = pluginRespData
		total := db.GetPluginsTotal(&field)
		data["total"] = total
		c.JSON(msg.SuccessResp(data))
		return
	} else {
		c.JSON(msg.ErrResp(msg.DealValidError(valid)))
		return
	}
}

//新增
func CreatePlugin(c *gin.Context) {
	plugin := db.Plugin{}
	err := c.BindJSON(&plugin)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	if db.ExistPluginByVulId(plugin.VulId){
		c.JSON(msg.ErrResp("漏洞编号已存在"))
		return
	} else {
		db.AddPlugin(plugin)
		c.JSON(msg.SuccessResp(plugin))
		return
	}
}

//修改
func UpdatePlugin(c *gin.Context) {
	plugin := db.Plugin{}
	err := c.BindJSON(&plugin)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	valid := validation.Validation{}
	valid.Min(plugin.Id, 1, "id").Message("ID必须大于0")
	valid.Required(plugin.Affects, "Affects").Message("Affects不能为空")

	if ! valid.HasErrors() {
		if db.ExistPluginByVulId(plugin.VulId){
			c.JSON(msg.ErrResp("漏洞编号已存在"))
			return
		} else {
			db.EditPlugin(plugin.Id, plugin)
			c.JSON(msg.SuccessResp(plugin))
			return
		}
	} else {
		c.JSON(msg.ErrResp(msg.DealValidError(valid)))
		return
	}
}

//删除
func DeletePlugin(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if ! valid.HasErrors() {
		if db.ExistPluginByID(id) {
			db.DeletePlugin(id)
			c.JSON(msg.SuccessResp("删除成功"))
			return
		} else {
			c.JSON(msg.ErrResp("record not found"))
			return
		}
	} else {
		c.JSON(msg.ErrResp(msg.DealValidError(valid)))
		return
	}
}

//运行单个plugin 不是从数据库提取数据，表单传数据
func RunPlugin(c *gin.Context) {
	run := RunSinglePluginSerializer{}
	err := c.BindJSON(&run)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	if run.Target != "" && run.JsonPoc != nil && run.Affects != "" {
		oreq, err := util.GenOriginalReq(run.Target)
		if err != nil {
			c.JSON(msg.ErrResp("原始请求生成失败"))
			return
		} else {
			poc, err := rule.ParseJsonPoc(run.JsonPoc)
			if err != nil {
				logging.GlobalLogger.Error("[plugins.go] fail to load plugins")
				c.JSON(msg.ErrResp("规则加载失败"))
			}
			currentPlugin := rule.Plugin{
				Affects:       run.Affects,
				JsonPoc:       poc,
			}
			item := &rule.ScanItem{Req: oreq, Plugin: &currentPlugin}
			result, err := rule.RunPoc(item)
			if err != nil {
				c.JSON(msg.ErrResp("规则运行失败：" + err.Error()))
				return
			}
			c.JSON(msg.SuccessResp(result))
			return
		}
	} else {
		c.JSON(msg.ErrResp("检测目标、规则类型均不可为空"))
		return
	}
}

//批量运行plugin 从数据库提取数据，表单传数据
//前端向后端传 "vul_id_list":["poc_db_1","poc_db_2"]
func RunPlugins(c *gin.Context) {
	runs := RunPluginsSerializer{}
	err := c.BindJSON(&runs)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	plugins, err := rule.LoadDbPlugin(runs.RunType, runs.VulIdList)

	switch runs.TargetType {
	case TargetUrl:
		url := runs.TargetType
		oreq, err := util.GenOriginalReq(url)
		if err != nil {
			logging.GlobalLogger.Error("[original request gen err ]", err)
			c.JSON(msg.ErrResp("原始请求生成失败"))
			return
		}
		rule.RunPlugins(oreq, plugins)
	case TargetUrlFile:
		//获取文件
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			logging.GlobalLogger.Error("[original request gen err ]", err)
			c.JSON(msg.ErrResp("url文件上传失败"))
			return
		}
		log.Print(header.Filename)
		//content, err := ioutil.ReadAll(file)
		var targets []string

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			val := scanner.Text()
			if val == "" {
				continue
			}
			targets = append(targets, val)
		}

		for _, url := range targets {
			oreq, err := util.GenOriginalReq(url)
			if err != nil {
				logging.GlobalLogger.Error("[original request gen err ]", err)
			}
			logging.GlobalLogger.Info("[start check url ]", url)
			rule.RunPlugins(oreq, plugins)
		}
	case TargetUrlRaw:
		//请求报文
	}
}





