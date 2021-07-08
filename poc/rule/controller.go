package rule

import (
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
	"sync"
)

const (
	AffectContent          = "text"
	AffectDirectory        = "directory"
	AffectURL              = "url"
	AffectAppendParameter  = "appendparam"
	AffectReplaceParameter = "replaceparam"
	AffectServer           = "server"
	AffectScript           = "script"
)

var ControllerPool = sync.Pool{}

func NewController() *PocController {
	if v := ControllerPool.Get(); v != nil {
		c := v.(*PocController)
		return c
	}
	return new(PocController)
}

func PutController(c *PocController) {
	c.Reset()
	ControllerPool.Put(c)
}

type PocController struct {
	Plugin		*Plugin
	Request		*RequestController
	CEL         *CelController
	Handles     []HandlerFunc          // 控制整个执行过程
	Index       int64         // 和middlefunc 配套
	abortIndex  int64         // 终止的index
	ScriptResult *util.ScanResult
	Debug		bool
	Keys        map[string]interface{}
}


type controllerContext interface {
	Next()
	Abort()
	IsAborted() bool
	GetString(key string) string
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	GetPoc() *Poc
	Groups(bool) (result bool, err error)
	Rules([]Rule, bool) (result bool, err error)
	GetPocName() string
	GetOriginalReq() *http.Request
	SetResult(result *util.ScanResult)
	IsDebug() bool
	// RegisterHandle(f HandlersChain)
}

func InitPocController(req *RequestController, plugin *Plugin, cel *CelController, handles []HandlerFunc) *PocController {
	controller := NewController()
	controller.Request = req
	controller.Plugin = plugin
	controller.CEL = cel
	controller.Handles = handles
	controller.Debug = false
	return controller
}

//	增加插件
func (controller *PocController) AddMiddle(handle HandlerFunc) {
	controller.Handles = append(controller.Handles, handle)
}

// 根据原始请求 + rule 生成并发起新的请求
func (controller *PocController) DoSingleRuleRequest(rule *Rule) (*proto.Response, error) {
	fastReq := controller.Request.Fast
	// fixReq : 根据规则对原始请求进行变形
	fixedFastReq := fasthttp.AcquireRequest()
	fastReq.CopyTo(fixedFastReq)
	curPath := string(fixedFastReq.URI().Path())

	affects := controller.Plugin.Affects

	switch affects {
	// param级
	case AffectAppendParameter, AffectReplaceParameter:
		for k, v := range rule.Headers {
			fixedFastReq.Header.Set(k, v)
		}
		return util.DoFasthttpRequest(fixedFastReq, rule.FollowRedirects)
	//	content级
	case AffectContent:
		return util.DoFasthttpRequest(fixedFastReq, rule.FollowRedirects)
	// dir级
	case AffectDirectory:
		// 目录级漏洞检测 判断是否以 "/"结尾
		if curPath != "" && strings.HasSuffix(curPath, "/") {
			// 去掉规则中的的"/" 再拼
			curPath = fmt.Sprint(curPath, strings.TrimPrefix(rule.Path, "/"))
		} else {
			curPath = fmt.Sprint(curPath, "/" ,strings.TrimPrefix(rule.Path, "/"))
		}
	// server级
	case AffectServer:
		curPath = rule.Path
	// url级
	case AffectURL:
		//curPath = curPath, strings.TrimPrefix(rule.Path, "/"))
	default:
	}
	// 兼容xray: 某些 POC 没有区分path和query
	curPath = strings.ReplaceAll(curPath, " ", "%20")
	curPath = strings.ReplaceAll(curPath, "+", "%20")

	fixedFastReq.URI().DisablePathNormalizing= true
	fixedFastReq.URI().Update(curPath)

	for k, v := range rule.Headers {
		fixedFastReq.Header.Set(k, v)
	}
	fixedFastReq.Header.SetMethod(rule.Method)

	// 处理multipart
	contentType := string(fixedFastReq.Header.ContentType())
	if strings.HasPrefix(strings.ToLower(contentType),"multipart/form-Data") && strings.Contains(rule.Body,"\n\n") {
		multipartBody, err := util.DealMultipart(contentType, rule.Body)
		if err != nil {
			return nil, err
		}
		fixedFastReq.SetBody([]byte(multipartBody))
	}else {
		fixedFastReq.SetBody([]byte(rule.Body))
	}
	return util.DoFasthttpRequest(fixedFastReq, rule.FollowRedirects)
}

// 单个规则运行
func (controller *PocController) SingleRule(rule *Rule, debug bool) (bool, error) {
	// 格式校验
	err := rule.Verify()
	if err != nil {
		return false, err
	}
	// 替换 set
	rule.ReplaceSet(controller.CEL.ParamMap)
	// 根据原始请求 + rule 生成并发起新的请求 http
	resp, err := controller.DoSingleRuleRequest(rule)
	if err != nil {
		return false, err
	}
	controller.CEL.ParamMap["response"] = resp
	// 匹配search规则
	if rule.Search != "" {
		controller.CEL.ParamMap = rule.ReplaceSearch(resp, controller.CEL.ParamMap)
	}
	// 如果当前rule验证失败，立即释放
	out, err := controller.CEL.Evaluate(rule.Expression)
	if err != nil {
		log.Error("[rule/controller.go:SingleRule cel evaluate error]", err)
		return false, err
	}
	if debug {
		controller.Request.Add(resp)
	} else {
		// 非debug模式下不记录 没有漏洞不记录请求链
		if !out {
			util.ResponsePut(resp)
			return false, nil
		} // 如果成功，记如请求链
		controller.Request.Add(resp)
	}
	return out, err
}

// 执行 rules
func (controller *PocController) Rules(rules []Rule, debug bool) (bool, error) {
	success := false
	for _, rule := range rules {
		singleRuleResult, err := controller.SingleRule(&rule, debug)
		if err != nil {
			log.Error("[rule/controller.go:Rules run error]" , err)
			return false, err
		}
		if !singleRuleResult {
			//如果false 直接跳出循环 返回
			success = false
			break
		}
		success = true
	}
	return success, nil
}

// 执行 groups
func (controller *PocController) Groups(debug bool) (bool, error) {
	groups := controller.Plugin.JsonPoc.Groups
	// groups 就是多个rules 任何一个rules成功 即返回成功
	for _, rules := range groups {
		rulesResult, err := controller.Rules(rules, debug)
		if err != nil || !rulesResult {
			continue
		}
		// groups中一个rules成功 即返回成功
		if rulesResult {
			return rulesResult, nil
		}
	}
	return false, nil
}

func (controller *PocController) Next() {

	for controller.Index < int64(len(controller.Handles)) {
		controller.Handles[controller.Index](controller)
		controller.Index++
	}
}


func (controller *PocController) IsAborted() bool {
	return controller.Index <= controller.abortIndex
}


func (controller *PocController) Abort() {
	controller.abortIndex = controller.Index + 1
}

func (controller *PocController) Reset() {
	controller.Handles = nil
	controller.Index = 0
	controller.abortIndex = 0
	controller.Plugin = nil
	controller.CEL.Reset()
	controller.Request.Reset()
	controller.ScriptResult = nil
	controller.Keys = nil
	controller.Debug = false
	return
}

func (controller *PocController) Set(key string, value interface{}) {
	if controller.Keys == nil {
		controller.Keys = make(map[string]interface{})
	}
	controller.Keys[key] = value
}

func (controller *PocController) Get(key string) (value interface{}, exists bool) {
	value, exists = controller.Keys[key]
	return
}

func (controller *PocController) GetString(key string) (s string) {
	if val, ok := controller.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (controller *PocController) GetPoc() *Poc {
	return controller.Plugin.JsonPoc
}

func (controller *PocController) GetPocName() string {
	return controller.Plugin.JsonPoc.Name
}

func (controller *PocController) GetOriginalReq() *http.Request {
	return controller.Request.Original
}

func (controller *PocController) SetResult(result *util.ScanResult){
	controller.ScriptResult = result
}

func (controller *PocController) IsDebug() bool {
	return controller.Debug
}
