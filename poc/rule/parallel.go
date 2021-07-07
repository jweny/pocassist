package rule

import (
	"errors"
	"github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	log "github.com/jweny/pocassist/pkg/logging"
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
	err := yaml.Unmarshal(jsonByte, poc)
	if poc.Name == "" {
		errMsg := "poc解析失败，poc名称不可为空"
		log.Error("rule/plugin.go:ParseJsonPoc Err", errMsg)
		return nil, errors.New(errMsg)
	}
	return poc, err
}

func ParseYamlPoc(yamlByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := yaml.Unmarshal(yamlByte, poc)
	if poc.Name == "" {
		errMsg := "poc parse error"
		log.Error("rule/plugin.go:ParseJsonPoc Err", errMsg)
		return nil, errors.New(errMsg)
	}
	return poc, err
}

// 限制并发
type ScanItem struct {
	OriginalReq *http.Request // 原始请求
	Plugin      *Plugin       // 检测插件
	Task        *db.Task      // 所属任务
}

func (item *ScanItem) Verify() error {
	errMsg := ""
	if item.Task == nil {
		errMsg = "task create fail"
		log.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	if item.OriginalReq == nil || item.Plugin == nil {
		errMsg = "not original request"
		log.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	if item.Plugin == nil {
		errMsg = "not plugin"
		log.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	return nil
}


//	从数据库 中加载 POC
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
			log.Error("[rule/parallel.go:LoadDbPlugin load multi err]", tx.Error)
			return nil, tx.Error
		}
	default:
		// 默认执行全部启用规则
		tx := db.GlobalDB.Where("enable = ?", 1).Find(&dbPluginList)
		if tx.Error != nil {
			log.Error("[rule/parallel.go:LoadDbPlugin load all err]", tx.Error)
			return nil, tx.Error
		}
	}
	log.Error("[rule/parallel.go:LoadDbPlugin load plugin number]", len(dbPluginList))

	for _, v := range dbPluginList {
		poc, err := ParseJsonPoc(v.JsonPoc)
		if err != nil {
			continue
		}
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


// 并发测试
func RunPlugins(plugins []Plugin, task *db.Task){
	var wg sync.WaitGroup

	// 插件并发数
	parallel := conf.GlobalConfig.PluginsConfig.Parallel

	p, _ := ants.NewPoolWithFunc(parallel, func(item interface{}) {
		RunPoc(item)
		wg.Done()
	})
	defer p.Release()

	for OriginalReq := range OriginalReqChannel{
		log.Info("[rule/parallel.go:RunPlugins start scan]", OriginalReq.URL.String())
		for i := range plugins {
			item := &ScanItem{OriginalReq, &plugins[i], task}
			wg.Add(1)
			p.Invoke(item)

		}
	}
	// todo 这里刷新task状态
	wg.Wait()
	db.DownTask(task.Id)
}