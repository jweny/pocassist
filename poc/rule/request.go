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
	OriginalQueryParams string
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
	rc.OriginalQueryParams = paramsString
	return nil
}

func (rc *RequestController) FixQueryParams(field string, payload string, affects string) (err error) {
	if rc.OriginalQueryParams == "" {
		err = rc.InitOriginalQueryParams()
	}
	if err != nil {
		return err
	}
	qs, err  := url.ParseQuery(rc.OriginalQueryParams)
	if err != nil {
		return err
	}
	var value string
	if vs, ok := qs[field]; ok {
		if len(vs) == 0 {
			value = ""
		} else {
			value = vs[0]
		}
		qs.Del(field)
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
	tmpQuery := qs.Encode()
	if tmpQuery != "" {
		tmpQuery += "&"
	}
	// 把`field`放在最后，供人工验证时判断
	tmpQuery += fmt.Sprintf("%v=%v", field, value)
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
		log.Error("util/requests.go:InitFast Err", err)
		return err
	}
	rc.Fast = fastReq
	return
}

func (rc *RequestController) InitData() (err error) {
	rc.Data, err = util.GetOriginalReqBody(rc.Original)
	if err != nil {
		log.Error("util/requests.go:InitData Err", err)
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


