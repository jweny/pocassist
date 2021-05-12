package rule

import (
	"encoding/json"
	"net/http"
	"pocassist/basic"
	"pocassist/database"
	"strings"
)

const (
	LoadSingle = "single"
	LoadAll = "all"
	LoadAffects = "affects"
	LoadMulti = "multi"
)

func ParseJsonPoc(jsonByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := json.Unmarshal(jsonByte, poc)
	return poc, err
}

// 按逗号切割 去除空格
func SplitToArray(conditions string) []string {
	array := strings.Split(conditions, ",")
	for index, value := range array {
		array[index] = strings.TrimSpace(value)
	}
	return array
}

//	从数据库 中加载 poc
func LoadDbPlugins(loadType string, conditions string) ([]database.Plugin, error) {
	var plugin database.Plugin
	var plugins []database.Plugin
	basic.GlobalLogger.Debug("[loadPoc type ]", loadType)
	basic.GlobalLogger.Debug("[conditions is ]", conditions)
	// todo 命令行里传json_str过来
	switch loadType {
	case LoadSingle:
		// 漏洞编号
		tx := database.GlobalDB.Where("vul_id = ? AND enable = ?", conditions, 1).First(&plugin)
		if tx.Error != nil {
			basic.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}
		plugins = append(plugins, plugin)

	case LoadAll:
		// 加载全部数据 无论是否启用
		tx := database.GlobalDB.Find(&plugins)
		if tx.Error != nil {
			basic.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}

	case LoadAffects:
		tx := database.GlobalDB.Where("affects = ? AND enable = ?", conditions, 1).Find(&plugins)
		if tx.Error != nil {
			basic.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}

	case LoadMulti:
		plugins := SplitToArray(conditions)
		tx := database.GlobalDB.Where("vul_id IN ? AND enable = ?", plugins, 1).Find(&plugins)
		if tx.Error != nil {
			basic.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}

	default:
		// 默认执行全部启用规则
		tx := database.GlobalDB.Where("enable = ?", 1).Find(&plugins)
		if tx.Error != nil {
			basic.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}
	}
	basic.GlobalLogger.Info("[plugins load number ]", len(plugins))
	return plugins, nil
}


// pluginsDB 转 plugin
func LoadPlugins(loadType string, conditions string) ([]Plugin, error) {
	var vuls []Plugin
	plugins, err := LoadDbPlugins(loadType, conditions)
	if err != nil {
		return nil, err
	}

	for _, v := range plugins {
		poc, err := ParseJsonPoc(v.JsonPoc)
		if err != nil {
			basic.GlobalLogger.Error("[plugins plugin load err ]", v.VulId)
			continue
		}
		rule := Plugin{
			VulId:         v.VulId,
			Affects:       v.Affects,
			JsonPoc:       poc,
			Enable:        v.Enable,
		}
		vuls = append(vuls, rule)
	}
	return vuls, nil
}

// 批量执行plugin
func RunPlugins(oreq *http.Request, rules []Plugin){
	for _, curRule := range rules {
		item := &ScanItem{oreq, &curRule}
		result, err := RunPoc(item)
		if err != nil {
			basic.GlobalLogger.Error("[plugins plugin run err ]", curRule.VulId)
		}
		basic.GlobalLogger.Info("[plugin result ]\n",
			" [vul_id] ", curRule.VulId,
			" [vul_name] ", curRule.JsonPoc.Name,
			" [vul_result] ", result)
	}
}