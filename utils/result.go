package utils

// 保存扫描结果
type ScanResult struct {
	Vulnerable bool     `json:"vulnerable"`// 是否存在漏洞
	Target     string   `json:"target"`// 漏洞url
	Output     string   `json:"output"`// 一些说明
	ReqMsg     []string `json:"req_msg"`// 请求列表
	RespMsg    []string `json:"resp_msg"`// 响应列表
}

type FormatString struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

// 没漏洞时返回的结果
var InVulnerableResult = ScanResult{
	Vulnerable: false,
}

// 有漏洞时返回的结果(http)
func VulnerableHttpResult(target string, output string, respList []*Response) *ScanResult {
	var reqMsg []string
	var respMsg []string
	defer ResponsesPut(respList)

	for _, v := range respList {
		reqMsg = append(reqMsg, v.ReqRaw)
		respMsg = append(respMsg, v.RespRaw)
	}
	return &ScanResult{
		Vulnerable: true,
		Target:     target,
		Output:     output,
		ReqMsg:     reqMsg,
		RespMsg:    respMsg,
	}
}

// 有漏洞时返回的结果(tcp/udp)
func VulnerableTcpOrUdpResult(target string, output string, payload []string, resp []string) *ScanResult {
	return &ScanResult{
		Vulnerable: true,
		Target:     target,
		Output:     output,
		ReqMsg:     payload,
		RespMsg:    resp,
	}
}
