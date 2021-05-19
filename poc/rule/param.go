package rule

import (
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"net/url"
	"strings"
)


type ReplaceHandler interface {
	Replace(value string, field string, controller *PocController)
}

type ReplaceGet struct {
}

func (r *ReplaceGet) Replace(value string, field string, controller *PocController) {
	reqQuery := ReplaceGetParam(controller.NewReq, value, field, controller.Affects)
	req := controller.NewReq
	curURL := fmt.Sprintf("%s://%s%s?%s", req.Url.Scheme, req.Url.Host, req.Url.Path, reqQuery)
	controller.FastReq.SetRequestURI(curURL)
	return
}

type ReplacePost struct {
}

func (r *ReplacePost) Replace(value string, field string, controller *PocController) {
	bodyString := ReplacePostParam(string(controller.reqData), value, field, controller.Affects)
	controller.FastReq.SetBodyString(bodyString)
	return
}

// 返回 get url
func ReplaceGetParam(originalReq *proto.Request, paramValue string, originalFields string, affects string) string {
	originalQuery := originalReq.Url.Query
	qs, err := url.ParseQuery(originalQuery)
	if err != nil {
		return ""
	}
	var value string
	if vs, ok := qs[originalFields]; ok {
		if len(vs) == 0 {
			value = ""
		} else {
			value = vs[0]
		}
		qs.Del(originalFields)
	} else {
		return originalQuery
	}

	if affects == AffectAppendParameter {
		value += paramValue
	} else {
		value = paramValue
	}
	tmpQuery := qs.Encode()
	if tmpQuery != "" {
		tmpQuery += "&"
	}
	// 需要把`field`放在最后，供二次验证时判断
	tmpQuery += fmt.Sprintf("%v=%v", originalFields, value)
	return tmpQuery
}

// 返回 post body
func ReplacePostParam(data string, paramValue string, originalFields string, affects string) string {
	qs, err := url.ParseQuery(data)
	if err != nil {
		return ""
	}
	value := qs.Get(originalFields)
	qs.Del(originalFields)
	if paramValue != "" {
		if affects == AffectAppendParameter {
			value += paramValue
		} else {
			value = paramValue
		}
	}
	tmpQuery := qs.Encode()
	if tmpQuery != "" {
		tmpQuery += "&"
	}
	// 需要把`field`放在最后，供二次验证时判断
	tmpQuery += fmt.Sprintf("%v=%v", originalFields, value)
	return tmpQuery
}

//	将 params 中的 {{}} 替换为自定义的set
func ParsePocParams(params []string, varMap map[string]interface{}) []string {
	// 把 param 替换为set参数
	for index, param := range params {
		if strings.Contains(param,"{{") && strings.Contains(param,"}}"){
			for setKey, setValue := range varMap {
				// 过滤掉 map
				_, isMap := setValue.(map[string]string)
				if isMap {
					continue
				}
				replaceStr := "{{"+setKey+"}}"
				if strings.Contains(param, replaceStr) {
					tmpStr := strings.ReplaceAll(param, replaceStr, fmt.Sprintf("%v", setValue))
					params[index] = tmpStr
				}
			}
		}
	}
	return params
}
