package rule

import (
	"encoding/json"
	"errors"
	"github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/panjf2000/ants/v2"
	"gopkg.in/yaml.v2"
	"net/http"
	"sync"
)

const (
	LoadMulti = "multi"
)

func ParseJsonPoc(jsonByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := json.Unmarshal(jsonByte, poc)
	if poc.Name == "" {
		return nil, errors.New("poc解析失败，poc名称不可为空")
	}
	return poc, err
}

func ParseYamlPoc(yamlByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := yaml.Unmarshal(yamlByte, poc)
	if poc.Name == "" {
		return nil, errors.New("poc解析失败，poc名称不可为空")
	}
	return poc, err
}

//	从数据库 中加载 poc
func LoadDbPlugin(lodeType string, array []string) ([]Plugin, error) {
	// 数据库数据
	var dbPluginList []db.Plugin
	// plugin对象
	var plugins []Plugin
	switch lodeType {
	case LoadMulti:
		// 多个
		tx := db.GlobalDB.Where("vul_id IN ? AND enable = ?", array, 1).Find(&dbPluginList)
		if tx.Error != nil {
			logging.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}
	default:
		// 默认执行全部启用规则
		tx := db.GlobalDB.Where("enable = ?", 1).Find(&dbPluginList)
		if tx.Error != nil {
			logging.GlobalLogger.Error("[db select err ]", tx.Error)
			return nil, tx.Error
		}
	}

	logging.GlobalLogger.Info("[dbPluginList load number ]", len(dbPluginList))
	logging.GlobalLogger.Info("[dbPluginList ]")
	for _, v := range dbPluginList {
		poc, err := ParseJsonPoc(v.JsonPoc)
		if err != nil {
			logging.GlobalLogger.Error("[plugins plugin load err ]", v.VulId)
			continue
		}
		logging.GlobalLogger.Info("[Plugin ]", poc.Name)
		plugin := Plugin{
			VulId:   v.VulId,
			Affects: v.Affects,
			JsonPoc: poc,
			Enable:  v.Enable,
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil

}

// 批量执行plugin
func RunPlugins(oreq *http.Request, plugins []Plugin, task *db.Task){
	// 并发限制
	var wg sync.WaitGroup
	parallel := conf.GlobalConfig.PluginsConfig.Parallel
	p, _ := ants.NewPoolWithFunc(parallel, func(item interface{}) {
		RunPoc(item)
		wg.Done()
	})
	defer p.Release()

	for i := range plugins {
		item := &ScanItem{oreq, &plugins[i], task}
		wg.Add(1)
		p.Invoke(item)
	}
	wg.Wait()
}