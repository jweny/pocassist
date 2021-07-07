package scripts

import (
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"regexp"
)

var jbossPattern1 = regexp.MustCompile(`org.jboss.invocation.InvocationException`)
var jbossPattern2 = regexp.MustCompile(`\$org.jboss.invocation.MarshalledValue`)

// JBossInvokerServletRemoteCodeExec jboss 远程执行
func JBossInvokerServletRemoteCodeExec(args *ScriptScanArgs) (*util.ScanResult, error) {
	// 定义报文列表
	var respList []*proto.Response

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.Header.SetMethod(fasthttp.MethodGet)

	for _, uri := range []string{
		"/invoker/EJBInvokerServlet/",
		"/invoker/JMXInvokerServlet/",
	} {
		rawUrl := ConstructUrl(args, uri)
		fastReq.SetRequestURI(rawUrl)
		resp, err := util.DoFasthttpRequest(fastReq, true)
		if err != nil {
			return nil, err
		}
		if resp.Status == 200 {
			if jbossPattern1.Match(resp.Body) && jbossPattern2.Match(resp.Body) {
				respList = append(respList, resp)
				return util.VulnerableHttpResult(rawUrl,"", respList),nil
			}
		}
		util.ResponsePut(resp)
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-jboss-invoker-servlet-rce", JBossInvokerServletRemoteCodeExec)
}
