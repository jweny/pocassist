package rule

import (
	"context"
	"errors"
	"golang.org/x/time/rate"
	"net/http"
	"pocassist/basic"
	"pocassist/scripts"
	"strconv"
	"sync"
	"time"
)

// 限制速率
var limiter *rate.Limiter

func InitRate() {
	maxQps := basic.GlobalConfig.HttpConfig.MaxQps
	parallel := basic.GlobalConfig.PluginsConfig.Parallel

	limit := rate.Every(time.Duration(maxQps) * time.Millisecond)
	// 第二个参数 和 并发加载的 plugin 数匹配
	limiter = rate.NewLimiter(limit, parallel)
}

func LimitWait() {
	limiter.Wait(context.Background())
}

// 限制并发
type ScanItem struct {
	Req *http.Request // 原始请求
	Vul *Plugin       // vul from db
}

var scanItemPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return new(ScanItem)
	},
}

func ScanItemPut(i interface{}) {
	item := i.(*ScanItem)
	item.Req = nil
	item.Vul = nil
	scanItemPool.Put(item)
	return
}

var Handles map[string][]HandlerFunc

func ExecExpressionHandle(controller *PocController) error {
	var result bool
	var err error
	if controller.poc.Groups != nil {
		result, err = controller.ParseGroupsRule()
	} else {
		result, err = controller.ParsePocRule()
	}
	if err != nil {
		return err
	}
	if result {
		controller.Abort()
	}

	basic.GlobalLogger.Debug("[plugin finish]", controller.poc.Name)
	return nil
}

func ExecScriptHandle(controller *PocController) error {
	scanFunc := scripts.GetScriptFunc(controller.poc.Name)
	if scanFunc == nil {
		return errors.New("未找到匹配的脚本，请检查脚本register方法中的第一个参数是否为规则名称")
	}

	var isHTTPS bool
	defaultPort := 80
	if controller.originalReq.URL.Scheme == "https" {
		isHTTPS = true
		defaultPort = 443
	}

	if controller.originalReq.URL.Port() != "" {
		port, err := strconv.ParseUint(controller.originalReq.URL.Port(), 10, 16)
		if err != nil {
			controller.Abort()
			return err
		}
		defaultPort = int(port)
	}

	args := &scripts.ScriptScanArgs{
		Host:    controller.originalReq.URL.Hostname(),
		Port:    uint16(defaultPort),
		IsHTTPS: isHTTPS,
	}

	result, err := scanFunc(args)
	if err != nil {
		basic.GlobalLogger.Error("[script scan failed ]", controller.vulId, " err:", err)
		return err
	}
	basic.GlobalLogger.Info("[script scan finished ]",
		" [vul_id] ", controller.vulId,
		" [script_func] ", scanFunc,
		" [vul_result] ", result)
	return nil
}

func InitHandles() {
	Handles = make(map[string][]HandlerFunc)
	Handles[AffectScript] = []HandlerFunc{ExecScriptHandle}
	Handles[AffectAppendParameter] = []HandlerFunc{ExecExpressionHandle}
	Handles[AffectReplaceParameter] = []HandlerFunc{ExecExpressionHandle}
}

func getHandles(affect string) []HandlerFunc {
	defaultHandlers := []HandlerFunc{ExecExpressionHandle}
	if handles, exists := Handles[affect]; exists {
		return handles
	}
	return defaultHandlers
}