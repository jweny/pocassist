package rule

import (
	"context"
	"errors"
	"github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/poc/scripts"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

// 限制速率
var limiter *rate.Limiter

func InitRate() {
	maxQps := conf.GlobalConfig.HttpConfig.MaxQps
	parallel := conf.GlobalConfig.PluginsConfig.Parallel

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
		logging.GlobalLogger.Info("[=== find vul===]\n",
			" [vul_id] ", controller.vulId,
			" [vul_name] ", controller.poc.Name)
		controller.Abort()
	}

	logging.GlobalLogger.Info("[=== not vul===]\n",
		" [vul_id] ", controller.vulId,
		" [vul_name] ", controller.poc.Name)
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
		logging.GlobalLogger.Error("[script scan failed ]", controller.vulId, " err:", err)
		return err
	}
	logging.GlobalLogger.Info("[script scan finished ]",
		" [vul_id] ", controller.vulId,
		" [script_func] ", scanFunc,
		" [vul_result] ", result)
	return nil
}

func Setup() {
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