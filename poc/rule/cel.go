package rule

import (
	"fmt"
	"github.com/google/cel-go/cel"
	cel2 "github.com/jweny/pocassist/pkg/cel"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/cel/reverse"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
)

type CelController struct {
	Env      *cel.Env               // cel env
	ParamMap map[string]interface{} // 注入到cel中的变量
}

//	初始化
func (cc *CelController) Init(poc *Poc) (err error) {
	//	1.生成cel env环境
	option := cel2.InitCelOptions()
	//	注入set定义的变量
	if poc.Set != nil {
		option.AddRuleSetOptions(poc.Set)
	}
	env, err := cel2.InitCelEnv(&option)
	if err != nil {
		log.Error("[rule/cel.go:Init init cel env error]", err)
		return err
	}
	cc.Env = env
	// 初始化变量列表
	cc.ParamMap = make(map[string]interface{})
	return nil
}

// 处理poc: set
func (cc *CelController) InitSet(poc *Poc, newReq *proto.Request) (err error) {
	// 如果没有set 就直接返回
	if len(poc.Set) == 0 {
		return
	}
	cc.ParamMap["request"] = newReq

	for _, setItem := range poc.Set {
		key := setItem.Key.(string)
		value := setItem.Value.(string)
		// 反连平台
		if value == "newReverse()" {
			cc.ParamMap[key] = reverse.NewReverse()
			continue
		}
		out, err := cel2.Evaluate(cc.Env, value, cc.ParamMap)
		if err != nil {
			return err
		}
		switch value := out.Value().(type) {
		// set value 无论是什么类型都先转成string
		case *proto.UrlType:
			cc.ParamMap[key] = util.UrlTypeToString(value)
		case int64:
			cc.ParamMap[key] = int(value)
		default:
			cc.ParamMap[key] = fmt.Sprintf("%v", out)
		}
	}
	return
}


// 计算cel表达式
func (cc *CelController) Evaluate(char string) (bool, error) {
	out, err := cel2.Evaluate(cc.Env, char, cc.ParamMap)
	if err != nil {
		log.Error("[rule/cel.go:Evaluate error]", err)
		return false, err
	}
	if fmt.Sprintf("%v", out) == "false"{
		return false, nil
	}
	return true, nil
}

func (cc *CelController) Reset(){
	cc.Env = nil
	cc.ParamMap = nil
}

