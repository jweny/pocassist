package scan

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/file"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"io/ioutil"
	"net/http"
	"path"
)

type scanSerializer struct {
	// 单个url
	Target  string   `json:"target" binding:"required"`
	Type    string   `json:"type" binding:"required,oneof=multi all"` // multi or all
	VulList []string `json:"vul_list"`
	Remarks string   `json:"remarks"`
}

type swaggerArray struct{
	VulList []string `json:"vul_list"`
}

// @Summary scan url
// @Tags Scan
// @Description 扫描单个url
// @accept json
// @Produce  json
// @Param scan body scanSerializer true "扫描参数"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/url [post]
func Url(c *gin.Context) {
	scan := scanSerializer{}
	err := c.ShouldBindJSON(&scan)
	if err != nil {
		c.JSON(msg.ErrResp("测试url不可为空，扫描类型为multi或all"))
		return
	}

	oreq, err := util.GenOriginalReq(scan.Target)
	if err != nil || oreq == nil {
		c.JSON(msg.ErrResp("原始请求生成失败"))
		return
	}

	// 插件列表
	plugins, err := rule.LoadDbPlugin(scan.Type, scan.VulList)
	if err != nil || plugins == nil{
		c.JSON(msg.ErrResp("poc插件加载失败" + err.Error()))
		return
	}
	token := c.Request.Header.Get("Authorization")
	claims, _ := util.ParseToken(token)

	// 创建任务
	task := db.Task{
		Operator: claims.Username,
		Remarks: scan.Remarks,
		Target:  scan.Target,
	}
	db.AddTask(&task)
	taskItem := &rule.TaskItem{
		OriginalReq: oreq,
		Plugins:     plugins,
		Task:        &task,
	}

	c.JSON(msg.SuccessResp("任务下发成功"))
	go rule.TaskProducer(taskItem)
	go rule.TaskConsumer()
	return
}

// @Summary scan raw
// @Tags Scan
// @Description 传文件：请求报文
// @Accept multipart/form-data
// @Param type formData string true "扫描类型：multi / all"
// @Param vul_list formData swaggerArray false "vul_id 列表"
// @Param remarks formData string false "备注"
// @Param target formData file true "file"
// @accept json
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/raw [post]
func Raw(c *gin.Context) {
	scanType := c.PostForm("type")
	vulList := c.PostFormArray("vul_list")
	remarks := c.PostForm("remarks")

	if scanType != "multi" && scanType != "all" {
		c.JSON(msg.ErrResp("扫描类型为multi或all"))
		return
	}

	target, err := c.FormFile("target")
	if err != nil {
		c.JSON(msg.ErrResp("文件上传失败"))
		return
	}
	// 存文件
	filePath := file.UploadTargetsPath(path.Ext(target.Filename))
	err = c.SaveUploadedFile(target, filePath)

	if err != nil || !file.Exists(filePath) {
		c.JSON(msg.ErrResp("文件保存失败"))
		return
	}

	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(msg.ErrResp("请求报文文件解析失败"))
		return
	}

	oreq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	if err != nil || oreq == nil {
		c.JSON(msg.ErrResp("生成原始请求失败"))
		return
	}
	if !oreq.URL.IsAbs() {
		scheme := "http"
		oreq.URL.Scheme = scheme
		oreq.URL.Host = oreq.Host
	}

	plugins, err := rule.LoadDbPlugin(scanType, vulList)
	if err != nil || plugins == nil {
		c.JSON(msg.ErrResp("插件加载失败" + err.Error()))
		return
	}

	oReqUrl := oreq.URL.String()

	token := c.Request.Header.Get("Authorization")
	claims, _ := util.ParseToken(token)

	task := db.Task{
		Operator: claims.Username,
		Remarks: remarks,
		Target:  oReqUrl,
	}
	db.AddTask(&task)
	taskItem := &rule.TaskItem{
		OriginalReq: oreq,
		Plugins:     plugins,
		Task:        &task,
	}

	c.JSON(msg.SuccessResp("任务下发成功"))
	go rule.TaskProducer(taskItem)
	go rule.TaskConsumer()
	return
}


// @Summary scan list
// @Tags Scan
// @Description 传文件：url列表
// @Accept multipart/form-data
// @Param type formData string true "扫描类型：multi / all"
// @Param vul_list formData swaggerArray false "vul_id 列表"
// @Param remarks formData string false "备注"
// @Param target formData file true "file"
// @Produce  json
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/list [post]
func List(c *gin.Context) {
	scanType := c.PostForm("type")
	vulList := c.PostFormArray("vul_list")
	remarks := c.PostForm("remarks")

	if scanType != "multi" && scanType != "all" {
		c.JSON(msg.ErrResp("扫描类型为multi或all"))
		return
	}

	target, err := c.FormFile("target")
	if err != nil {
		c.JSON(msg.ErrResp("文件上传失败"))
		return
	}
	// 存文件
	filePath := file.UploadTargetsPath(path.Ext(target.Filename))
	err = c.SaveUploadedFile(target, filePath)

	if err != nil || !file.Exists(filePath) {
		c.JSON(msg.ErrResp("文件保存失败"))
		return
	}

	// 加载poc
	plugins, err := rule.LoadDbPlugin(scanType, vulList)
	if err != nil{
		c.JSON(msg.ErrResp("插件加载失败" + err.Error()))
		return
	}
	if len(plugins) == 0 {
		c.JSON(msg.ErrResp("插件加载失败" + err.Error()))
		return
	}
	targets := file.ReadingLines(filePath)

	token := c.Request.Header.Get("Authorization")
	claims, _ := util.ParseToken(token)

	var oReqList []*http.Request
	var taskList []*db.Task

	for _, url := range targets {
		oreq, err := util.GenOriginalReq(url)
		if err != nil {
			continue
		}
		task := db.Task{
			Operator: claims.Username,
			Remarks: remarks,
			Target:  url,
		}
		db.AddTask(&task)

		oReqList = append(oReqList, oreq)
		taskList = append(taskList, &task)
	}

	if len(oReqList) == 0 || len(taskList) ==0 {
		c.JSON(msg.ErrResp("url列表加载失败"))
		return
	}
	c.JSON(msg.SuccessResp("任务下发成功"))

	for index, oreq := range oReqList {
		taskItem := &rule.TaskItem{
			OriginalReq: oreq,
			Plugins:     plugins,
			Task:        taskList[index],
		}
		go rule.TaskProducer(taskItem)
		go rule.TaskConsumer()
	}
	return
}

