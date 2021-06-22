package rule

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetOriginalReqBody(originalReq *http.Request) ([]byte, error){
	var data []byte
	if originalReq.Body != nil && originalReq.Body != http.NoBody {
		data, err := ioutil.ReadAll(originalReq.Body)
		if err != nil {
			return nil, err
		}
		originalReq.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return data, nil
}

func WriteVulResult(scanItem *ScanItem, output string, result *util.ScanResult){
	vulnerable := false
	if result != nil {
		vulnerable = result.Vulnerable
	}

	logging.GlobalLogger.Info(
		"[plugin scan result]",
		" [exist vul] ", vulnerable,
		" [plugin_id] ", scanItem.Plugin.VulId,
		" [plugin_name] ", scanItem.Plugin.JsonPoc.Name,
		" [output] ", output,
		" [detail] ", result)
	detail, _:= json.Marshal(result)

	res := db.Result{
		Detail:   detail,
		PluginId:  scanItem.Plugin.VulId,
		PluginName:	scanItem.Plugin.JsonPoc.Name,
		Vul:	vulnerable,
		TaskId: scanItem.Task.Id,
	}
	db.AddResult(res)
}

func WriteTaskError(errMsg string, taskId int) {
	logging.GlobalLogger.Error("[RunPoc]" + errMsg)
	// 修改task状态
	db.ErrorTask(taskId)
}

// 执行单个poc
func RunPoc(inter interface{}) (*util.ScanResult, error) {
	scanItem := inter.(*ScanItem)
	originalReq := scanItem.Req
	plugin := scanItem.Plugin
	task := scanItem.Task

	if originalReq == nil || plugin == nil {
		WriteTaskError("no request or no plugin", task.Id)
		return nil, errors.New("no request or no plugin")
	}

	data, err := GetOriginalReqBody(originalReq)
	if err != nil {
		WriteTaskError("original request body get err", task.Id)
		return nil, err
	}

	handles := getHandles(plugin.Affects)
	logging.GlobalLogger.Debug("[plugin running ]" , plugin.VulId, " [affects] ", plugin.Affects, " [name] ", plugin.JsonPoc.Name)

	// =========================
	env, err := GenCelEnv(plugin.JsonPoc)
	if err != nil {
		WriteTaskError("cel env gen err", task.Id)
		return nil, err
	}
	newReq, err := InitNewReq(originalReq)
	if err != nil {
		WriteTaskError("new request init err", task.Id)
		return nil, err
	}
	varMap, err := ParsePocSet(plugin.JsonPoc, env, newReq)
	if err != nil {
		util.RequestPut(newReq)
		WriteTaskError("poc set parse err", task.Id)
		return nil, err
	}
	// ===========================

	// 影响为参数类型  需要替换参数
	if plugin.Affects == AffectReplaceParameter || plugin.Affects == AffectAppendParameter {
		var originalGetParamFields url.Values
		var replaceHandler ReplaceHandler
		var err error
		if originalReq.Method == "GET" {
			originalGetParamFields, err = url.ParseQuery(originalReq.URL.RawQuery)
			if err != nil {
				WriteTaskError("originalReqGET url parse err", task.Id)
				return nil, err
			}
			replaceHandler = &ReplaceGet{}
		} else if originalReq.Method == "POST" {
			originalGetParamFields, err = url.ParseQuery(string(data))
			if err != nil {
				WriteTaskError("originalReqPost url parse err", task.Id)
				return nil, err
			}
			replaceHandler = &ReplacePost{}
		}

		for field := range originalGetParamFields {
			for _, value := range plugin.JsonPoc.Params {
				// 限速
				LimitWait()
				logging.GlobalLogger.Debug("[current param]", value)

				controller := InitPocController(originalReq, plugin.JsonPoc, plugin.Affects, data)
				controller.celEnv = env
				controller.varMap = varMap
				controller.Handles = handles
				controller.NewReq = newReq

				replaceHandler.Replace(value, field, controller)
				err := controller.Next()
				if err != nil {
					return nil, err
				}

				if controller.IsAborted() {
					result := util.VulnerableHttpResult(controller.originalReq.URL.String(),value, controller.respList)
					WriteVulResult(scanItem, value, result)
					// 修改task状态
					db.DownTask(scanItem.Task.Id)
					util.RequestPut(newReq)
					controller.Reset()
					return result, nil
				}
				controller.Reset()
			}
		}
		util.RequestPut(newReq)
		WriteVulResult(scanItem, "", nil)
	} else {
		// 影响为其他类型
		// 限速
		LimitWait()
		controller := InitPocController(originalReq, plugin.JsonPoc, plugin.Affects, data)
		controller.celEnv = env
		controller.varMap = varMap
		controller.Handles = handles
		controller.pluginId = plugin.VulId

		err = controller.Next()
		if err != nil {
			return nil, err
		}

		if controller.IsAborted() {
			result := util.VulnerableHttpResult(controller.originalReq.URL.String(),"", controller.respList)
			WriteVulResult(scanItem, "", result)
			util.RequestPut(newReq)
			// 修改task状态
			db.DownTask(scanItem.Task.Id)
			controller.Reset()
			return result, nil
		}
		WriteVulResult(scanItem, "", nil)
		controller.Reset()
		util.RequestPut(newReq)
	}
	// 修改task状态
	db.DownTask(scanItem.Task.Id)
	return &util.InVulnerableResult, nil
}
