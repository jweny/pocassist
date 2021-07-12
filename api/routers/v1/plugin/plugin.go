package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/unknwon/com"
	"gopkg.in/yaml.v2"
	"gorm.io/datatypes"
	"io/ioutil"
)

type Serializer struct {
	// 返回给前端的字段
	DespName		string		   `json:"desp_name"`
	Id       		int            `json:"id"`
	VulId           string         `json:"vul_id"`
	Affects         string         `json:"affects"`
	JsonPoc         datatypes.JSON `json:"json_poc"`
	Enable          bool           `json:"enable"`
	Description   	int            `json:"description"`
}

type RunSerializer struct {
	// 运行单个
	Target			string		   `json:"target" binding:"required"`
	VulId			string		   `json:"vul_id"`
	Affects         string         `json:"affects" binding:"required"`
	JsonPoc         datatypes.JSON `json:"json_poc"`
}

// 这个结构体没用 只是为了定义 swigger 前端的参数
type RunSwigger struct {
	Target			string		   `json:"target"`
	VulId			string		   `json:"vul_id"`
	Affects         string         `gorm:"column:affects" json:"affects"`
	JsonPoc         rule.Plugin `gorm:"column:json_poc" json:"json_poc"`
}

type DownloadSerializer struct {
	// 下载 yaml
	JsonPoc         datatypes.JSON `json:"json_poc"`
}

// 这个结构体没用 只是为了定义 swigger 前端的参数
type DownloadSwigger struct {
	JsonPoc         rule.Plugin `gorm:"column:json_poc" json:"json_poc"`
}

// @Summary plugin detail
// @Tags Plugin
// @Description 详情
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [get]
func Detail(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface {}
	if id < 0 {
		c.JSON(msg.ErrResp("ID必须大于0"))
		return
	}
	if db.ExistPluginByID(id) {
		data = db.GetPlugin(id)
		c.JSON(msg.SuccessResp(data))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}

// @Summary plugin list
// @Tags Plugin
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.PluginSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/ [get]
func Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := db.PluginSearchField{Search: "", EnableField:-1, AffectsField:"",}

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
	}
	if arg := c.Query("affectsField"); arg != "" {
		field.AffectsField = arg
	}

	plugins := db.GetPlugins(page, pageSize, &field)

	var pluginRespData []Serializer

	for _, plugin := range plugins {
		var despName string
		if plugin.Vulnerability != nil {
			despName = plugin.Vulnerability.NameZh
		} else {
			despName = ""
		}

		pluginRespData = append(pluginRespData, Serializer{
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
}

// @Summary plugin add
// @Tags Plugin
// @Description 新增
// @Produce  json
// @Security token
// @Param plugin body rule.Plugin true "plugin"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/ [post]
func Add(c *gin.Context) {
	plugin := db.Plugin{}
	err := c.ShouldBindJSON(&plugin)
	if err != nil {
		c.JSON(msg.ErrResp("参数不合法"))
		return
	}
	// 漏洞编号自动生成
	vulId, err := db.GenPluginVulId()
	if err != nil {
		c.JSON(msg.ErrResp("漏洞编号生成失败"))
		return
	}
	plugin.VulId = vulId
	db.AddPlugin(plugin)
	c.JSON(msg.SuccessResp(plugin))
	return
}

// @Summary plugin update
// @Tags Plugin
// @Description 更新
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Param plugin body rule.Plugin true "plugin"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [put]
func Update(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	plugin := db.Plugin{}
	err := c.ShouldBindJSON(&plugin)
	if err != nil {
		c.JSON(msg.ErrResp("参数不合法"))
		return
	}
	if db.ExistPluginByID(id) {
		db.EditPlugin(id, plugin)
		c.JSON(msg.SuccessResp(plugin))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}

// @Summary plugin delete
// @Tags Plugin
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [delete]
func Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	if db.ExistPluginByID(id) {
		db.DeletePlugin(id)
		c.JSON(msg.SuccessResp("删除成功"))
		return
	} else {
		c.JSON(msg.ErrResp("record not found"))
		return
	}
}

// @Summary plugin run
// @Tags Plugin
// @Description 运行
// @Produce  json
// @Security token
// @Param run body RunSwigger false "run"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/run/ [post]
func Run(c *gin.Context) {
	run := RunSerializer{}
	err := c.ShouldBindJSON(&run)
	if err != nil {
		c.JSON(msg.ErrResp("漏洞编号以poc-开头 测试url和漏洞类型不可为空"))
		return
	}

	oreq, err := util.GenOriginalReq(run.Target)
	if err != nil {
		c.JSON(msg.ErrResp("原始请求生成失败"))
		return
	}
	verify := util.VerifyTargetConnection(oreq)
	if !verify {
		c.JSON(msg.ErrResp("测试目标连通性测试不通过"))
		return
	}
	poc, err := rule.ParseJsonPoc(run.JsonPoc)
	if err != nil {
		log.Error("[plugins.go run] fail to load plugins", err)
		c.JSON(msg.ErrResp("规则加载失败"))
		return
	}

	task := db.Task{
		Remarks:  "single poc",
		Target:   run.Target,
	}
	db.AddTask(&task)

	currentPlugin := rule.Plugin{
		Affects:       run.Affects,
		JsonPoc:       poc,
		VulId: 		   run.VulId,
	}

	item := &rule.ScanItem{oreq, &currentPlugin, &task}

	result, err := rule.RunPoc(item, true)
	if err != nil {
		db.ErrorTask(task.Id)
		c.JSON(msg.ErrResp("规则运行失败：" + err.Error()))
		return
	}
	db.DownTask(task.Id)
	c.JSON(msg.SuccessResp(result))
	return
}

// @Summary download yaml
// @Tags Plugin
// @Description 下载yaml
// @Security token
// @Param run body DownloadSwigger true "json_poc"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/download/ [post]
func DownloadYaml(c *gin.Context) {
	download := DownloadSerializer{}
	err := c.ShouldBindJSON(&download)
	if err != nil {
		c.JSON(msg.ErrResp("规则格式不正确"))
		return
	}
	poc, err := rule.ParseJsonPoc(download.JsonPoc)
	if err != nil {
		log.Error("[plugins.go download] fail to load plugins", err)
		c.JSON(msg.ErrResp("规则解析失败"))
		return
	}
	content, err :=yaml.Marshal(poc)
	if err != nil {
		log.Error("[plugins.go download] fail to marshal yaml", err)
		c.JSON(msg.ErrResp("yaml生成失败"))
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", poc.Name))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Writer.Write(content)
}

// @Summary upload yaml
// @Tags Plugin
// @Description 上传yaml
// @Accept multipart/form-data
// @Param yaml formData file true "file"
// @accept json
// @Security token
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/upload/ [post]
func UploadYaml(c *gin.Context) {
	file, _, err := c.Request.FormFile("yaml")
	if err != nil {
		c.JSON(msg.ErrResp("文件上传失败"))
		return
	}
	// 获取yaml内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(msg.ErrResp("文件读取失败，请检查后重试"))
		return
	}
	poc, err := rule.ParseYamlPoc(content)
	if err != nil {
		c.JSON(msg.ErrResp("yaml解析失败，请检查后重试"))
		return
	}
	// todo slice to map
	toMap := TempPoc{
		Params: poc.Params,
		Name:   poc.Name,
		Set:    SliceToMap(poc.Set),
		Rules:  poc.Rules,
		Groups: poc.Groups,
		Detail: rule.Detail{},
	}
	data := make(map[string]interface{})
	data["json_poc"] = toMap
	c.JSON(msg.SuccessResp(data))
}

type TempPoc struct {
	Params	[]string	 	  `json:"params"`
	Name   string             `json:"name"`
	Set    map[string]string  `json:"set"`
	Rules  []rule.Rule        `json:"rules"`
	Groups  map[string][]rule.Rule `json:"groups"`
	Detail rule.Detail             `json:"detail"`
}
func SliceToMap(slice yaml.MapSlice) map[string]string {
	m := make(map[string]string)
	for _,v := range slice{
		m[v.Key.(string)] = v.Value.(string)
	}
	return m
}
