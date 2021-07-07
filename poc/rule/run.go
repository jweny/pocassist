package rule

import (
	"encoding/json"
	"github.com/jweny/pocassist/pkg/db"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
)

func WriteTaskResult(scanItem *ScanItem, result *util.ScanResult) {
	log.Info("[rule/poc.go:WriteTaskResult scan finish]", scanItem.OriginalReq.URL.String(), scanItem.Plugin.JsonPoc.Name, result)
	vulnerable := false
	if result != nil {
		vulnerable = result.Vulnerable
	}
	detail, _:= json.Marshal(result)

	res := db.Result{
		Detail:   detail,
		PluginId:  scanItem.Plugin.VulId,
		PluginName:	scanItem.Plugin.JsonPoc.Name,
		Vul:	vulnerable,
		TaskId: scanItem.Task.Id,
	}
	db.AddResult(res)
	return
}

// 执行单个poc
func RunPoc(inter interface{}) (result *util.ScanResult, err error) {
	scanItem := inter.(*ScanItem)
	err = scanItem.Verify()
	if err != nil {
		log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
		return nil, err
	}
	log.Info("[rule/poc.go:RunPoc current plugin]", scanItem.Plugin.JsonPoc.Name)

	var requestController RequestController
	var celController CelController

	err = requestController.Init(scanItem.OriginalReq)
	if err != nil {
		log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
		return nil, err
	}

	handles := getHandles(scanItem.Plugin.Affects)

	err = celController.Init(scanItem.Plugin.JsonPoc)
	if err != nil {
		log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
		return nil, err
	}

	err = celController.InitSet(scanItem.Plugin.JsonPoc, requestController.New)
	if err != nil {
		util.RequestPut(requestController.New)
		log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
		return nil, err
	}
	switch scanItem.Plugin.Affects {
	// 影响为参数类型
	case AffectAppendParameter, AffectReplaceParameter:
		{
			err := requestController.InitOriginalQueryParams()
			if err != nil {
				log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
				return nil, err
			}
			controller := InitPocController(&requestController, scanItem.Plugin, &celController, handles)
			paramPayloadList := scanItem.Plugin.JsonPoc.Params

			for field := range requestController.OriginalQueryParams {
				for _, payload := range paramPayloadList {
					// 限速
					LimitWait()
					log.Info("[rule/poc.go:RunPoc param payload]", payload)
					err = requestController.FixQueryParams(field, payload, controller.Plugin.Affects)
					if err != nil {
						log.Error("[rule/poc.go:WriteTaskError scan error] ", err)
						continue
					}
					controller.Next()
					if controller.IsAborted() {
						// 存在漏洞
						result = util.VulnerableHttpResult(controller.GetOriginalReq().URL.String(), payload, controller.Request.Raw)
						WriteTaskResult(scanItem, result)
						PutController(controller)
						return result, nil
					}
				}
			}
			// 没漏洞
			WriteTaskResult(scanItem, &util.InVulnerableResult)
			return &util.InVulnerableResult, nil
		}
	case AffectDirectory, AffectServer, AffectURL:
		{
			// todo 报错 刷任务状态
			LimitWait()
			controller := InitPocController(&requestController, scanItem.Plugin, &celController, handles)
			controller.Next()
			if controller.IsAborted() {
				// 存在漏洞
				result = util.VulnerableHttpResult(controller.GetOriginalReq().URL.String(), "", controller.Request.Raw)
				WriteTaskResult(scanItem, result)
				PutController(controller)
				return result, nil
			} else {
				// 没漏洞
				WriteTaskResult(scanItem, &util.InVulnerableResult)
				PutController(controller)
				return &util.InVulnerableResult, nil
			}
		}
	case AffectScript:
		{
			LimitWait()
			controller := InitPocController(&requestController, scanItem.Plugin, &celController, handles)
			controller.Next()
			if controller.IsAborted() && controller.ScriptResult != nil {
				// 存在漏洞 保存结果
				result = &util.ScanResult{
					Vulnerable: controller.ScriptResult.Vulnerable,
					Target:     controller.ScriptResult.Target,
					Output:     controller.ScriptResult.Output,
					ReqMsg:     controller.ScriptResult.ReqMsg,
					RespMsg:    controller.ScriptResult.RespMsg,
				}
				WriteTaskResult(scanItem, controller.ScriptResult)
				PutController(controller)
				return result, nil
			} else {
				// 没漏洞
				WriteTaskResult(scanItem, &util.InVulnerableResult)
				PutController(controller)
				return &util.InVulnerableResult, nil
			}
		}
	}
	// 默认返回没有漏洞
	return &util.InVulnerableResult, nil
}