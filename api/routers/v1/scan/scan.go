package scan

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"log"
)

type UrlItem struct {
	// 批量运行
	Target    string   `json:"target"`
	LoadType  string   `json:"run_type"`    // multi or all
	VulIdList []string `json:"vul_id_list"` //前端向后端传 "vul_id_list":["poc_db_1","poc_db_2"]
}

// 单个url
func Url(c *gin.Context) {
	item := UrlItem{}
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	plugins, err := rule.LoadDbPlugin(item.LoadType, item.VulIdList)
	oreq, err := util.GenOriginalReq(item.Target)
	if err != nil {
		logging.GlobalLogger.Error("[original request gen err ]", err)
		c.JSON(msg.ErrResp("原始请求生成失败"))
		return
	}
	// todo 加载config
	ch := make(chan util.ScanResult, 100)
	rule.RunPlugins(oreq, plugins)
}

// 加载文件 批量扫描
func File(c *gin.Context) {
	//获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logging.GlobalLogger.Error("[original request gen err ]", err)
		c.JSON(msg.ErrResp("url文件上传失败"))
		return
	}
	log.Print(header.Filename)
	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		if val == "" {
			continue
		}
		targets = append(targets, val)
	}

	item := UrlItem{}
	err = c.BindJSON(&item)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}

	for _, url := range targets {
		plugins, err := rule.LoadDbPlugin(item.LoadType, item.VulIdList)
		oreq, err := util.GenOriginalReq(url)
		if err != nil {
			logging.GlobalLogger.Error("[original request gen err ]", err)
			c.JSON(msg.ErrResp("原始请求生成失败"))
			return
		}
		rule.RunPlugins(oreq, plugins)
	}
}
