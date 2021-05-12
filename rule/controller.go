package rule

import (
	"errors"
	"fmt"
	"net/http"
	"pocassist/basic"
	"pocassist/utils"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/valyala/fasthttp"
)

type HandlerFunc func(*PocController) error

var ControllerPool = sync.Pool{
	New: func() interface{} {
		return new(PocController)
	},
}

type PocController struct {
	vulId	   string
	originalReq *http.Request 				// 原始请求  --> 初始条件
	poc        *Poc 						// 加载的poc --> 初始条件
	NewReq     *utils.Request				// 生成的新请求
	celEnv     *cel.Env						// cel env
	varMap     map[string]interface{} 		// 注入到cel中的变量
	fastReq    *fasthttp.Request      		// fasthttp 请求
	affects    string        				// 检测类型
	reqData    []byte        				// 请求的内容
	Handles    []HandlerFunc 				// 控制整个执行过程
	Index      int64         				// 和Handles配套
	abortIndex int64        	 			// 终止的index
	respList    []*utils.Response			// 记录请求和响应
}

func (controller *PocController) Next() error {
	for controller.Index < int64(len(controller.Handles)) {
		err := controller.Handles[controller.Index](controller)
		if err != nil {
			return err
		}
		controller.Index++
	}
	return nil
}

func (controller *PocController) IsAborted() bool {
	return controller.Index <= controller.abortIndex
}


func (controller *PocController) Abort() {
	controller.abortIndex = controller.Index + 1
}

func (controller *PocController) Reset() {
	fasthttp.ReleaseRequest(controller.fastReq)
	utils.ResponsesPut(controller.respList)
	controller.Handles = nil
	controller.Index = 0
	controller.abortIndex = 0
	controller.varMap = nil
	controller.reqData = nil
	controller.poc = nil
	controller.celEnv = nil
	controller.NewReq = nil
	controller.vulId = ""
	return
}

func InitPocController(originalReq *http.Request, p *Poc, affects string, data []byte) *PocController {
	controller := ControllerPool.Get().(*PocController)
	controller.originalReq = originalReq
	controller.poc = p
	controller.affects = affects
	controller.fastReq = fasthttp.AcquireRequest()
	controller.reqData = data
	utils.CopyRequest(originalReq, controller.fastReq, data)
	return controller
}

//	增加插件
func (controller *PocController) AddMiddle(handle HandlerFunc) {
	controller.Handles = append(controller.Handles, handle)
}

//	初始化表达式
func (controller *PocController) GenCelEnv() error {
	//	初始化表达式
	option := utils.InitCelOptions()
	//	注入set定义的变量
	if controller.poc.Set != nil {
		option.AddRuleSetOptions(controller.poc.Set)
	}
	//	生成cel环境
	env, err := utils.InitCelEnv(&option)
	if err != nil {
		basic.GlobalLogger.Error("[plugin cel env init ]" , err)
		return err
	}
	controller.celEnv = env
	return nil
}

//	初始化表达式
func GenCelEnv(poc *Poc) (env *cel.Env, err error) {
	//	初始化表达式
	option := utils.InitCelOptions()
	//	注入set定义的变量
	if poc.Set != nil {
		option.AddRuleSetOptions(poc.Set)
	}
	//	生成cel环境
	env, err = utils.InitCelEnv(&option)
	if err != nil {
		return
	}

	return
}


// 排序
func SortMapKeys(m map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// 处理poc: {{}} 替换为自定义的set
func ParsePocSingleRule(rule *Rule, varMap map[string]interface{}) *Rule {
	for setKey, setValue := range varMap {
		// 过滤掉 map
		_, isMap := setValue.(map[string]string)
		if isMap {
			continue
		}
		value := fmt.Sprintf("%v", setValue)
		// 替换请求头中的 自定义字段
		for headerKey, headerValue := range rule.Headers {
			rule.Headers[headerKey] = strings.ReplaceAll(headerValue, "{{"+setKey+"}}", value)
		}
		// 替换请求路径中的 自定义字段
		rule.Path = strings.ReplaceAll(strings.TrimSpace(rule.Path), "{{"+setKey+"}}", value)
		// 替换body的 自定义字段
		rule.Body = strings.ReplaceAll(strings.TrimSpace(rule.Body), "{{"+setKey+"}}", value)
	}
	return rule
}

// 实现 search
func doSearch(re string, body string) map[string]string {
	r, err := regexp.Compile(re)
	if err != nil {
		return nil
	}
	result := r.FindStringSubmatch(body)
	names := r.SubexpNames()
	if len(result) > 1 && len(names) > 1 {
		paramsMap := make(map[string]string)
		for i, name := range names {
			if i > 0 && i <= len(result) {
				paramsMap[name] = result[i]
			}
		}
		return paramsMap
	}
	return nil
}

// 处理poc: search
func ParsePocRuleSearch(rule *Rule, resp *utils.Response, varMap map[string]interface{}) map[string]interface{} {
	result := doSearch(strings.TrimSpace(rule.Search), string(resp.Body))
	if result != nil && len(result) > 0 { // 正则匹配成功
		for k, v := range result {
			varMap[k] = v
		}
	}
	return varMap
}

// 处理poc: set
func ParsePocSet(poc *Poc, env *cel.Env, newReq *utils.Request) (varMap map[string]interface{}, err error) {
	varMap = make(map[string]interface{})

	// 如果没有set 就直接返回
	if len(poc.Set) == 0 {
		return
	}
	varMap["request"] = newReq
	//	获取所有Set key
	setKeys := SortMapKeys(poc.Set)
	// 处理set 先排序解析除了payload以外，其他的自定义变量
	for _, k := range setKeys {
		setValue := poc.Set[k]
		if k != "payload" {
			if setValue == "newReverse()" {
				varMap[k] = utils.NewReverse()
				continue
			}
			out, err := utils.Evaluate(env, setValue, varMap)
			if err != nil {
				continue
			}
			switch value := out.Value().(type) {
			// set value 无论是什么类型都先转成string
			case *utils.UrlType:
				varMap[k] = utils.UrlTypeToString(value)
			case int64:
				varMap[k] = int(value)
			default:
				varMap[k] = fmt.Sprintf("%v", out)
			}
		}
	}
	// 最后处理payload
	if poc.Set["payload"] != "" {
		out, err := utils.Evaluate(env, poc.Set["payload"], varMap)
		if err != nil {
			return varMap, err
		}
		varMap["payload"] = fmt.Sprintf("%v", out)
	}
	return
}

// 执行 rules
func (controller *PocController) ParsePocRule() (bool, error) {
	success := false

	for _, rule := range controller.poc.Rules {
		// 限制rule中的path必须以"/"开头
		if strings.HasPrefix(rule.Path, "/") == false {
			return false, errors.New("poc rule path must startWith \"/\"")
		}
		// 替换 set
		completedRule := ParsePocSingleRule(&rule, controller.varMap)
		// 根据原始请求 + rule 生成并发起新的请求 http

		resp, err := DoSingleRuleRequest(controller, completedRule)
		if err != nil {
			basic.GlobalLogger.Error("[plugin http err ]" , err)
			return false, err
		}
		controller.varMap["response"] = resp
		// 匹配search规则
		if rule.Search != "" {
			controller.varMap = ParsePocRuleSearch(&rule, resp, controller.varMap)
		}
		out, err := utils.Evaluate(controller.celEnv, rule.Expression, controller.varMap)
		if err != nil {
			basic.GlobalLogger.Error("[plugin cel evaluate ]" , err)
			return false, err
		}
		if fmt.Sprintf("%v", out) == "false" {
			utils.ResponsePut(resp)
			//如果false不继续执行后续rule
			success = false
			// 其中任何一次失败，都将直接跳出循环
			break
		}

		basic.GlobalLogger.Info("============")
		basic.GlobalLogger.Info("req:", resp.ReqRaw)
		basic.GlobalLogger.Info("resp:", resp.RespRaw)
		controller.respList = append(controller.respList, resp)
		basic.GlobalLogger.Info("============")
		success = true
	}
	return success, nil
}

// 执行 groups
func (controller *PocController) ParseGroupsRule() (bool, error) {

	success := false

	for _, rules := range controller.poc.Groups {
		for _, rule := range rules {
			// 限制rule中的path必须以"/"开头
			if strings.HasPrefix(rule.Path, "/") == false {
				return false, errors.New("poc rule path must startWith \"/\"")
			}
			completedRule := ParsePocSingleRule(&rule, controller.varMap)
			// 请求
			resp, err := DoSingleRuleRequest(controller, completedRule)
			if err != nil {
				basic.GlobalLogger.Error("[plugin http err ]" , err)
				return false, err
			}
			controller.varMap["response"] = resp
			// 匹配search规则
			if rule.Search != "" {
				controller.varMap = ParsePocRuleSearch(&rule, resp, controller.varMap)
			}
			out, err := utils.Evaluate(controller.celEnv, rule.Expression, controller.varMap)
			if err != nil {
				basic.GlobalLogger.Error("[plugin cel evaluate ]" , err)
				return false, err
			}
			if fmt.Sprintf("%v", out) == "false" {
				utils.ResponsePut(resp)
				success = false
				// 其中任何一次失败，都将直接跳出循环
				break
			}
			f := utils.ReqFormat{
				Req: controller.fastReq,
			}
			resp.ReqRaw = f.FormatContent()

			controller.respList = append(controller.respList, resp)
			success = true
		}
		// groups中一个rules成功 即返回成功
		if success {
			return success, nil
		}
	}
	return success, nil
}
