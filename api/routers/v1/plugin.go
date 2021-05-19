package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/unknwon/com"
	"gorm.io/datatypes"
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

type RunPluginSerializer struct {
	// 运行
	Target			string		`json:"target"`
	Affects         string         `gorm:"column:affects" json:"affects"`
	JsonPoc         datatypes.JSON `gorm:"column:json_poc" json:"json_poc"`
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
		if db.ExistPluginByID(plugin.Id){
			db.EditPlugin(plugin.Id, plugin)
			c.JSON(msg.SuccessResp(plugin))
		} else {
			c.JSON(msg.ErrResp("record not found"))
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

//运行
func RunPlugin(c *gin.Context) {
	run := RunPluginSerializer{}
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
				logging.GlobalLogger.Error("[plugins plugin load err ]",)
				c.JSON(msg.ErrResp("规则加载失败"))
			}
			currentPlugin := rule.Plugin{
				Affects:       run.Affects,
				JsonPoc:       poc,
			}
			item := &rule.ScanItem{Req: oreq, Vul: &currentPlugin}
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




