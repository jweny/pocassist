package util

import (
	"github.com/jweny/pocassist/pkg/cel/proto"
)

// 保存扫描结果
type ScanResult struct {
	Vulnerable bool     `json:"vulnerable"`// 是否存在漏洞
	Target     string   `json:"target"`// 漏洞url
	Output     string   `json:"output"`// 一些说明
	ReqMsg     []string `json:"req_msg"`// 请求列表
	RespMsg    []string `json:"resp_msg"`// 响应列表
}

// 没漏洞时返回的结果
var InVulnerableResult = ScanResult{
	Vulnerable: false,
}

// debug没漏洞返回的结果(http)
func DebugVulnerableHttpResult(target string, output string, respList []*proto.Response) *ScanResult {
	var reqMsg []string
	var respMsg []string
	defer ResponsesPut(respList)

	for _, v := range respList {
		reqMsg = append(reqMsg, v.ReqRaw)
		respMsg = append(respMsg, v.RespRaw)
	}
	return &ScanResult{
		Vulnerable: false,
		Target:     target,
		Output:     output,
		ReqMsg:     reqMsg,
		RespMsg:    respMsg,
	}
}


// 有漏洞时返回的结果(http)
func VulnerableHttpResult(target string, output string, respList []*proto.Response) *ScanResult {
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
