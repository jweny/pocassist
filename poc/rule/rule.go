package rule

import (
	"errors"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	log "github.com/jweny/pocassist/pkg/logging"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
)

// 单个规则
type Rule struct {
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	Search          string            `json:"search"`
	FollowRedirects bool              `json:"follow_redirects"`
	Expression      string            `json:"expression"`
}

type Detail struct {
	Author      string   `json:"author"`
	Links       []string `json:"links"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
}

// Rules 和 Groups 只能存在一个
type Poc struct {
	Params	[]string	 	  `json:"params"`
	Name   string             `json:"name"`
	Set    yaml.MapSlice 	  `json:"set"`
	Rules  []Rule             `json:"rules"`
	Groups  map[string][]Rule `json:"groups"`
	Detail Detail             `json:"detail"`
}

type Plugin struct {
	VulId   string `gorm:"column:vul_id"` // 漏洞编号
	Affects string `gorm:"column:affects"`   // 影响类型  dir/server/param/url/content
	JsonPoc *Poc   `gorm:"column:json_poc"`  // json规则
	Enable  bool   `gorm:"column:enable"`    // 是否启用
}

// set
func (rule *Rule) ReplaceSet (varMap map[string]interface{}) {
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
}

// search
func (rule *Rule) ReplaceSearch(resp *proto.Response, varMap map[string]interface{}) map[string]interface{} {
	result := doSearch(strings.TrimSpace(rule.Search), string(resp.Body))
	if result != nil && len(result) > 0 { // 正则匹配成功
		for k, v := range result {
			varMap[k] = v
		}
	}
	return varMap
}

// 校验rule格式
func (rule *Rule) Verify () error {
	// 限制rule中的path必须以"/"开头
	if strings.HasPrefix(rule.Path, "/") == false {
		errorMsg := "POC rule path must startWith \"/\""
		log.Error("rule/rule.go:Verify error]" , errorMsg)
		return errors.New(errorMsg)
	}
	return nil
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