package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/conf"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 限制请求速率
var limiter *rate.Limiter

func InitRate() {
	msQps := conf.GlobalConfig.HttpConfig.MaxQps  / 10
	limit := rate.Every(100 * time.Millisecond)
	limiter = rate.NewLimiter(limit, msQps)
}

func LimitWait() {
	limiter.Wait(context.Background())
}


type clientDoer interface {
	// 不跟随重定向
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, t time.Duration) error
	// 跟随重定向
	DoRedirects(req *fasthttp.Request, resp *fasthttp.Response, maxRedirectsCount int) error
}

var (
	fasthttpClient   clientDoer
)

var (
	requestPool sync.Pool = sync.Pool{
		New: func() interface{} {
			return new(proto.Request)
		},
	}

	responsePool sync.Pool = sync.Pool{
		New: func() interface{} {
			return new(proto.Response)
		},
	}

	formatPool sync.Pool = sync.Pool{
		New: func() interface{} {
			return new(FormatString)
		},
	}
)

type FormatString struct {
	Raw string `json:"raw"`
}

type ReqFormat struct {
	Req *fasthttp.Request
}

type RespFormat struct {
	Resp *fasthttp.Response
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

// dump 请求报文
func (r *ReqFormat) FormatContent() string {
	reqRaw := formatPool.Get().(*FormatString)
	defer formatPut(reqRaw)
	req := r.Req
	if req == nil {
		return ""
	}
	// fasthttp 请求打印的第一行长这样
	// GET http://jweny.top/ HTTP/1.1
	// 处理下
	tmpList := strings.SplitN(string(req.Header.Header()), "\r\n",2)

	reqURI := req.URI().RequestURI()
	protocol := string(req.Header.Protocol())
	body := ""
	if len(r.Req.Body()) > 0 {
		body = string(r.Req.Body())
	}

	line1 := fmt.Sprintf("%s %s %s\r\n", valueOrDefault(string(req.Header.Method()), "GET"),
		reqURI, protocol)
	line2 := ""
	// 避免打印的 Host 头重复
	if !strings.Contains(tmpList[1], "Host:"){
		line2 = fmt.Sprintf("%s: %s\r\n", "Host", string(req.Host()))
	}
	requestRaw := line1 + line2 + tmpList[1] + body
	return requestRaw
}

// dump 响应报文
func (r *RespFormat) FormatContent() string {
	respRaw := formatPool.Get().(*FormatString)
	defer formatPut(respRaw)
	header := r.Resp.Header.String()
	body := ""
	if len(r.Resp.Body()) > 0 {
		body = string(r.Resp.Body())
	}
	responseRaw := header + body
	return responseRaw
}

func formatPut(f *FormatString) {
	if f == nil {
		return
	}
	f.Raw = ""
	formatPool.Put(f)
}

func RequestGet() *proto.Request {
	return requestPool.Get().(*proto.Request)
}

func RequestPut(r *proto.Request) {
	if r == nil {
		return
	}
	r.Reset()
	requestPool.Put(r)
}

func RespGet() *proto.Response {
	return responsePool.Get().(*proto.Response)
}

func ResponsePut(resp *proto.Response) {
	if resp == nil {
		return
	}
	resp.Reset()
	responsePool.Put(resp)
}

func ResponsesPut(responses []*proto.Response) {
	for _, item := range responses {
		ResponsePut(item)
	}
}

func ParseUrl(u *url.URL) *proto.UrlType {
	nu := &proto.UrlType{}
	nu.Scheme = u.Scheme
	nu.Domain = u.Hostname()
	nu.Host = u.Host
	nu.Port = u.Port()
	nu.Path = u.EscapedPath()
	nu.Query = u.RawQuery
	nu.Fragment = u.Fragment
	return nu
}

func ParseFasthttpResponse(originalResp *fasthttp.Response, req *fasthttp.Request) (*proto.Response, error) {
	resp := RespGet()
	header := make(map[string]string)
	resp.Status = int32(originalResp.StatusCode())
	u, err := url.Parse(req.URI().String())
	if err != nil {
		log.Error("util/requests.go:ParseFasthttpResponse url parse error", req.URI().String(), err)
		return nil, err
	}
	resp.Url = ParseUrl(u)

	headerContent := originalResp.Header.String()
	headers := strings.Split(headerContent, "\r\n")
	for _, v := range headers {
		// 修复bug: 限制切割次数
		values := strings.SplitN(v, ":", 2)
		if len(values) != 2 {
			continue
		}
		// 修复bug 所有响应头 key 均转为小写（与xray兼容）
		k := strings.ToLower(values[0])
		// 修复bug 所有响应头 value去除左边空格
		v := strings.TrimLeft(values[1]," ")
		// 修复bug 处理响应头 中多个相同的key 产生的覆盖问题
		if header[k] != "" {
			header[k] += v
		} else {
			header[k] = v
		}
	}
	resp.Headers = header
	resp.ContentType = string(originalResp.Header.Peek("Content-Type"))

	resp.Body = make([]byte, len(originalResp.Body()))
	copy(resp.Body, originalResp.Body())
	return resp, nil
}

func DoFasthttpRequest(req *fasthttp.Request, redirect bool) (*proto.Response, error) {
	LimitWait()
	defer fasthttp.ReleaseRequest(req)
	bodyLen := len(req.Body())
	if bodyLen > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(bodyLen))
		if string(req.Header.Peek("Content-Type")) == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	var err error

	if redirect {
		// 跟随重定向 最大跳转数从conf中加载
		maxRedirects := conf.GlobalConfig.HttpConfig.MaxRedirect
		err = fasthttpClient.DoRedirects(req, resp, maxRedirects)
	} else {
		// 不跟随重定向
		timeout := conf.GlobalConfig.HttpConfig.HttpTimeout
		err = fasthttpClient.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
	}
	if err != nil {
		log.Error("util/requests.go:DoFasthttpRequest fasthttp client doRequest error", string(req.RequestURI()),err)
		return nil, err
	}

	// 处理响应 body: gzip deflate 解包
	fixBody, err := UnzipResponseBody(resp)
	if err != nil {
		log.Error("util/requests.go:DoFasthttpRequest fasthttp client dealResponseBody error", string(req.RequestURI()),err)
		return nil, err
	}
	resp.SetBody(fixBody)
	curResp, err := ParseFasthttpResponse(resp, req)
	// 添加请求和响应报文
	if err != nil {
		return nil, err
	}

	f := RespFormat{
		Resp: resp,
	}
	curResp.RespRaw = f.FormatContent()

	reqf := ReqFormat{
		Req: req,
	}
	curResp.ReqRaw = reqf.FormatContent()
	return curResp, err
}

func UrlTypeToString(u *proto.UrlType) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Scheme != "" || u.Host != "" {
		if u.Host != "" || u.Path != "" {
			buf.WriteString("//")
		}
		if h := u.Host; h != "" {
			buf.WriteString(u.Host)
		}
	}
	path := u.Path
	if path != "" && path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)

	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}

func CopyRequest(req *http.Request, dstRequest *fasthttp.Request, data []byte) error {

	curURL := req.URL.String()
	dstRequest.SetRequestURI(curURL)
	dstRequest.Header.SetMethod(req.Method)

	for name, values := range req.Header {
		// Loop over all values for the name.
		for index, value := range values {
			if index > 0 {
				dstRequest.Header.Add(name, value)
			} else {
				dstRequest.Header.Set(name, value)
			}
		}
	}
	dstRequest.SetBodyRaw(data)
	return nil
}

// UnzipResponseBody 返回解压缩的 Body : 目前支持 identity gzip deflate
func UnzipResponseBody(response *fasthttp.Response) ([]byte, error) {
	contentEncoding := strings.ToLower(string(response.Header.Peek("Content-Encoding")))
	var body []byte
	var err error
	switch contentEncoding {
	case "", "none", "identity":
		body, err = response.Body(), nil
	case "gzip":
		body, err = response.BodyGunzip()
	case "deflate":
		body, err = response.BodyInflate()
	default:
		body, err = []byte{}, fmt.Errorf("unsupported Content-Encoding: %v", contentEncoding)
	}
	return body, err
}

func GenOriginalReq(url string) (*http.Request, error) {
	// 生成原始请求
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
	} else {
		url = "http://" + url
	}
	originalReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("util/requests.go:GenOriginalReq original request gen error", url, err)
		return nil, err
	}
	originalReq.Header.Set("Host", originalReq.Host)
	originalReq.Header.Set("Accept-Encoding", "gzip, deflate")
	originalReq.Header.Set("Accept","*/*")
	originalReq.Header.Set("User-Agent", conf.GlobalConfig.HttpConfig.Headers.UserAgent)
	originalReq.Header.Set("Accept-Language","en")
	originalReq.Header.Set("Connection","close")

	return originalReq, nil
}

func VerifyTargetConnection(originalReq *http.Request) bool {
	fastReq := fasthttp.AcquireRequest()
	fastResp := fasthttp.AcquireResponse()
	oriData, err := GetOriginalReqBody(originalReq)
	if err != nil {
		return false
	}
	err = CopyRequest(originalReq, fastReq, oriData)
	if err != nil {
		return false
	}
	timeout := conf.GlobalConfig.HttpConfig.HttpTimeout
	// 检测原始请求
	err = fasthttpClient.DoTimeout(fastReq, fastResp, time.Duration(timeout)*time.Second)
	if err != nil {
		// 检测原始请求 + index.php
		uri := string(fastReq.RequestURI())
		if uri != "" && strings.HasSuffix(uri, "/") {
			uri = fmt.Sprint(uri, "index.php")
		} else {
			uri = fmt.Sprint(uri, "/index.php")
		}
		fastReq.SetRequestURI(uri)
		err = fasthttpClient.DoTimeout(fastReq, fastResp, time.Duration(timeout)*time.Second)
		if err != nil {
			return false
		}
	}
	return true
}

func GetOriginalReqBody(originalReq *http.Request) ([]byte, error){
	var data []byte
	if originalReq.Body != nil && originalReq.Body != http.NoBody {
		data, err := ioutil.ReadAll(originalReq.Body)
		if err != nil {
			return nil, err
		}
		originalReq.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return data, nil
}


func DealMultipart(contentType string, ruleBody string) (result string, err error) {
	errMsg := ""
	// 处理multipart的/n
	re := regexp.MustCompile(`(?m)multipart\/form-Data; boundary=(.*)`)
	match := re.FindStringSubmatch(contentType)
	if len(match) != 2 {
		errMsg = "no boundary in content-type"
		//logging.GlobalLogger.Error("util/requests.go:DealMultipart Err", errMsg)
		return "", errors.New(errMsg)
	}
	boundary := "--" + match[1]
	multiPartContent := ""

	// 处理rule
	multiFile := strings.Split(ruleBody, boundary)
	if len(multiFile) == 0 {
		errMsg = "ruleBody.Body multi content format err"
		//logging.GlobalLogger.Error("util/requests.go:DealMultipart Err", errMsg)
		return multiPartContent, errors.New(errMsg)
	}

	for _, singleFile := range multiFile {
		//	处理单个文件
		//	文件头和文件响应
		spliteTmp := strings.Split(singleFile,"\n\n")
		if len(spliteTmp) == 2 {
			fileHeader := spliteTmp[0]
			fileBody := spliteTmp[1]
			fileHeader = strings.Replace(fileHeader,"\n","\r\n",-1)
			multiPartContent += boundary + fileHeader + "\r\n\r\n" + strings.TrimRight(fileBody ,"\n") + "\r\n"
		}
	}
	multiPartContent += boundary + "--" + "\r\n"
	return multiPartContent, nil
}

