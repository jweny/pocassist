package rule

import (
	"bytes"
	"errors"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 执行单个poc
func RunPoc(inter interface{}) (*util.ScanResult, error) {
	scanItem := inter.(*ScanItem)
	originalReq := scanItem.Req
	plugin := scanItem.Plugin

	if originalReq == nil || plugin == nil {
		return nil, errors.New("no request or no plugin")
	}

	var data []byte
	if originalReq.Body != nil && originalReq.Body != http.NoBody {
		data, err := ioutil.ReadAll(originalReq.Body)
		if err != nil {
			logging.GlobalLogger.Error("[plugin originalReq data read err ]", plugin.VulId)
			return nil, err
		}
		originalReq.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	handles := getHandles(plugin.Affects)
	logging.GlobalLogger.Debug("[plugin running ]" , plugin.VulId, " [affects] ", plugin.Affects, " [name] ", plugin.JsonPoc.Name)
	// 影响为参数类型
	if plugin.Affects == AffectReplaceParameter || plugin.Affects == AffectAppendParameter {
		var originalGetParamFields url.Values
		var replaceHandler ReplaceHandler
		var err error
		if originalReq.Method == "GET" {
			originalGetParamFields, err = url.ParseQuery(originalReq.URL.RawQuery)
			if err != nil {
				logging.GlobalLogger.Error("[plugin originalReqGET url parse err ]", err)
				return nil, err
			}
			replaceHandler = &ReplaceGet{}
		} else if originalReq.Method == "POST" {
			originalGetParamFields, err = url.ParseQuery(string(data))
			if err != nil {
				logging.GlobalLogger.Error("[plugin originalReqPost url parse err ]", err)
				return nil, err
			}
			replaceHandler = &ReplacePost{}
		}

		env, err := GenCelEnv(plugin.JsonPoc)
		if err != nil {
			logging.GlobalLogger.Error("[plugin cel env gen err ]" , plugin.VulId)
			return nil, err
		}
		newReq, err := InitNewReq(originalReq)
		if err != nil {
			logging.GlobalLogger.Error("[plugin new request init err ]" , plugin.VulId)
			return nil, err
		}
		varMap, err := ParsePocSet(plugin.JsonPoc, env, newReq)
		if err != nil {
			util.RequestPut(newReq)
			logging.GlobalLogger.Error("[plugin poc set parse err ]", plugin.VulId, err)
			return nil, err
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
					controller.Reset()
					util.RequestPut(newReq)
					logging.GlobalLogger.Info("[plugin result ]\n",
						" [plugin_id] ", plugin.VulId,
						" [plugin_name] ", plugin.JsonPoc.Name,
						" [param] ", value)

					return util.VulnerableHttpResult(controller.originalReq.URL.String(),"", controller.respList), nil
				}
				controller.Reset()
				util.RequestPut(newReq)
			}
		}

	} else {
		// 其他类型
		env, err := GenCelEnv(plugin.JsonPoc)
		if err != nil {
			logging.GlobalLogger.Error("[plugin cel env gen err ]" , plugin.VulId)
			return nil, err
		}
		newReq, err := InitNewReq(originalReq)
		if err != nil {
			logging.GlobalLogger.Error("[plugin new request init err ]" , plugin.VulId)
			return nil, err
		}
		varMap, err := ParsePocSet(plugin.JsonPoc, env, newReq)
		if err != nil {
			util.RequestPut(newReq)
			logging.GlobalLogger.Error("plugin poc set parse err ]", plugin.VulId, err)
			return nil, err
		}
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
			logging.GlobalLogger.Info("[plugin scan result ]\n",
				" [plugin_id] ", plugin.VulId,
				" [plugin_name] ", plugin.JsonPoc.Name,
				" [plugin_result] ", result)
			return result, nil
		}
		controller.Reset()
		util.RequestPut(newReq)
	}
	return &util.InVulnerableResult, nil
}
