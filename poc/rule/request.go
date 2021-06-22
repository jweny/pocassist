package rule

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func InitNewReq(originalReq *http.Request) (req *proto.Request, err error) {
	req = util.RequestGet()

	req.Method = originalReq.Method
	req.Url = util.ParseUrl(originalReq.URL)

	header := make(map[string]string)
	for k := range originalReq.Header {
		header[k] = originalReq.Header.Get(k)
	}
	req.Headers = header
	req.ContentType = originalReq.Header.Get("Content-Type")
	if originalReq.Body == nil || originalReq.Body == http.NoBody {
	} else {
		data, err := ioutil.ReadAll(originalReq.Body)
		if err != nil {
			return req, err
		}
		req.Body = data
		originalReq.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
	return
}


// 根据原始请求 + rule 生成并发起新的请求
func DoSingleRuleRequest(controller *PocController, rule *Rule) (*proto.Response, error) {
	affects := controller.Affects

	oldReq := controller.FastReq
	httpRequest := fasthttp.AcquireRequest()

	oldReq.CopyTo(httpRequest)

	curPath := string(httpRequest.URI().Path())

	switch affects {
	// param级
	case AffectAppendParameter, AffectReplaceParameter:
		for k, v := range rule.Headers {
			httpRequest.Header.Set(k, v)
		}
		return util.DoFasthttpRequest(httpRequest, rule.FollowRedirects)
	//	content级
	case AffectContent:
		return util.DoFasthttpRequest(httpRequest, rule.FollowRedirects)
	// dir级
	case AffectDirectory:
		// 目录级漏洞检测 判断是否以 "/"结尾
		if curPath != "" && strings.HasSuffix(curPath, "/") {
			// 去掉规则中的的"/" 再拼
			curPath = fmt.Sprint(curPath, strings.TrimPrefix(rule.Path, "/"))
		} else {
			//return nil, errors.New("affects Dir, but target url not a Dir")
			curPath = fmt.Sprint(curPath, '/' ,strings.TrimPrefix(rule.Path, "/"))
		}
	// server级
	case AffectServer:
		curPath = rule.Path
	// url级
	case AffectURL:
		//curPath = curPath, strings.TrimPrefix(rule.Path, "/"))
	default:
	}
	// 兼容xray: 某些 poc 没有区分path和query
	curPath = strings.ReplaceAll(curPath, " ", "%20")
	curPath = strings.ReplaceAll(curPath, "+", "%20")

	httpRequest.URI().DisablePathNormalizing= true
	httpRequest.URI().Update(curPath)

	for k, v := range rule.Headers {
		httpRequest.Header.Set(k, v)
	}
	httpRequest.Header.SetMethod(rule.Method)

	// 处理multipart
	contentType := string(httpRequest.Header.ContentType())
	if strings.HasPrefix(strings.ToLower(contentType),"multipart/form-data") && strings.Contains(rule.Body,"\n\n") {
		multipartBody, err := DealMultipart(contentType, rule.Body)
		if err != nil {
			return nil, err
		}
		httpRequest.SetBody([]byte(multipartBody))
	}else {
		httpRequest.SetBody([]byte(rule.Body))
	}
	return util.DoFasthttpRequest(httpRequest, rule.FollowRedirects)
}

func DealMultipart(contentType string, ruleBody string) (result string, err error) {
	// 处理multipart的/n
	re := regexp.MustCompile(`(?m)multipart\/form-data; boundary=(.*)`)
	match := re.FindStringSubmatch(contentType)
	if len(match) != 2 {
		return "", errors.New("no boundary in content-type")
	}
	boundary := "--" + match[1]
	multiPartContent := ""

	// 处理rule
	multiFile := strings.Split(ruleBody, boundary)
	if len(multiFile) == 0 {
		return multiPartContent, errors.New("ruleBody.Body multi content format err")
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


