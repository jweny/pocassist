package rule

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	log "github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// poc运行期间的各类请求
type RequestController struct{
	// 原始请求
	Original *http.Request
	// 经过变形的新请求
	New      *proto.Request
	// 真正发起的请求：转为 fasthttp
	Fast     *fasthttp.Request
	// post data
	Data     []byte
	// 记录请求和响应报文列表
	Raw      []*proto.Response
	// 原始请求的参数
	OriginalQueryParams url.Values
}

func (rc *RequestController) Init(original *http.Request) (err error) {
	rc.InitOriginal(original)
	err = rc.InitData()
	if err != nil {
		return err
	}
	err = rc.InitNew()
	if err != nil {
		return err
	}
	err = rc.InitFast()
	return err
}

func (rc *RequestController) InitOriginal(original *http.Request) {
	rc.Original = original
}

func (rc *RequestController) InitOriginalQueryParams() error {
	var paramsString string
	method := rc.Original.Method
	if method == "GET"{
		paramsString = rc.Original.URL.RawQuery
	}
	if method == "POST"{
		paramsString = string(rc.Data)
	}
	params, err := url.ParseQuery(paramsString)
	if err != nil {
		return err
	}
	rc.OriginalQueryParams = params
	return nil
}

func (rc *RequestController) FixQueryParams(field string, payload string, affects string) (err error) {
	if rc.OriginalQueryParams == nil {
		err = rc.InitOriginalQueryParams()
	}
	var value string
	if vs, ok := rc.OriginalQueryParams[field]; ok {
		if len(vs) == 0 {
			value = ""
		} else {
			value = vs[0]
		}
		rc.OriginalQueryParams.Del(field)
	} else {

		return errors.New("param payload fix err, field" + field + " not found")
	}

	if affects == AffectAppendParameter {
		value += payload
	} else if affects == AffectReplaceParameter{
		value = payload
	} else {
		return  errors.New("affects " + affects + " not support")
	}
	tmpQuery := rc.OriginalQueryParams.Encode()
	if tmpQuery != "" {
		tmpQuery += "&"
	}
	// 把`field`放在最后，供人工验证时判断
	tmpQuery += fmt.Sprintf("%v=%v", rc.OriginalQueryParams, value)
	method := rc.Original.Method
	if method == "GET"{
		currentUrl := fmt.Sprintf("%s://%s%s?%s", rc.New.Url.Scheme, rc.New.Url.Host, rc.New.Url.Path, tmpQuery)
		rc.Fast.SetRequestURI(currentUrl)
	} else {
		rc.Fast.SetBodyString(tmpQuery)
	}
	return nil
}

func (rc *RequestController) InitNew() (err error) {
	rc.New = util.RequestGet()
	rc.New.Method = rc.Original.Method
	rc.New.Url = util.ParseUrl(rc.Original.URL)

	header := make(map[string]string)
	for k := range rc.Original.Header {
		header[k] = rc.Original.Header.Get(k)
	}
	rc.New.Headers = header
	rc.New.ContentType = rc.Original.Header.Get("Content-Type")
	if rc.Original.Body == nil || rc.Original.Body == http.NoBody {
	} else {
		data, err := ioutil.ReadAll(rc.Original.Body)
		if err != nil {
			log.Error("rule/requests.go:InitNew gen request data error", err)
			return err
		}
		rc.New.Body = data
		rc.Original.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return
}

// 原始请求转为fasthttp
func (rc *RequestController) InitFast() (err error){
	fastReq := fasthttp.AcquireRequest()
	err = util.CopyRequest(rc.Original, fastReq, rc.Data)
	if err != nil {
		//logging.GlobalLogger.Error("util/requests.go:InitFast Err", err)
		return err
	}
	rc.Fast = fastReq
	return
}

func (rc *RequestController) InitData() (err error) {
	rc.Data, err = GetOriginalReqBody(rc.Original)
	if err != nil {
		//logging.GlobalLogger.Error("util/requests.go:InitData Err", err)
	}
	return err
}

func (rc *RequestController) Add(resp *proto.Response) {
	rc.Raw = append(rc.Raw, resp)
}

func (rc *RequestController) Reset() {
	fasthttp.ReleaseRequest(rc.Fast)
	util.RequestPut(rc.New)
	util.ResponsesPut(rc.Raw)
	rc.Data = nil
	rc.New = nil
	rc.Original = nil
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

