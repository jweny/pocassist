package utils

import (
	"fmt"
	"pocassist/basic"
	"testing"
)

func TestReverseCheck(t *testing.T) {
	InitFastHttpClient("")
	basic.InitLog(false,"")
	r := NewReverse()
	fmt.Println(r.Flag)
	fmt.Println(r.Domain)
	fmt.Println(r.Url)

	fmt.Println(ReverseCheck(r, 10))
}

//package utils
//
//import (
//	"bytes"
//	"fmt"
//	"github.com/valyala/fasthttp"
//	"net/url"
//	"pocassist/basic"
//	"time"
//)
//
//func NewReverse() *Reverse {
//	reverseUrl := basic.GlobalConfig.Reverse.Http.Url
//	reverseDomain := basic.GlobalConfig.Reverse.Dns.Domain
//	flag := RandLowLetterNumber(16)
//	webLogUrl, _ := url.Parse(fmt.Sprintf("%s/check-%s", reverseUrl, flag))
//	basic.GlobalLogger.Debug("[new reverse flag ]", )
//	return &Reverse{
//		Flag: flag,
//		// web log
//		Url: ParseUrl(webLogUrl),
//		// dns log
//		Domain:             flag + reverseDomain,
//		Ip:                 "",
//		IsDomainNameServer: false,
//	}
//}
//
//func ReverseCheck(r *Reverse, timeout int64) bool {
//	reverseUrl := basic.GlobalConfig.Reverse.Http.Url
//	reverseDomain := basic.GlobalConfig.Reverse.Dns.Domain
//	// 延迟 x 秒获取结果
//	time.Sleep(time.Second * time.Duration(timeout))
//	// check dns
//	verifyUrl := fmt.Sprintf("%s/%s-verify.php?verify&rmd=%s", reverseDomain, "dns", r.Flag)
//	if GetReverseResp(verifyUrl){
//		return true
//	} else {
//	//	check web
//		verifyUrl := fmt.Sprintf("%s/%s-verify.php?verify&rmd=%s", reverseUrl, "vul", r.Flag)
//		if GetReverseResp(verifyUrl){
//			return true
//		}
//	}
//	return false
//}
//
//func GetReverseResp(verifyUrl string) bool {
//	existsStr := []byte("Vulnerabilities exist")
//
//	fastReq := fasthttp.AcquireRequest()
//	defer fasthttp.ReleaseRequest(fastReq)
//	fastReq.SetRequestURI(verifyUrl)
//	fastReq.Header.SetMethod(fasthttp.MethodGet)
//
//	resp, err := DoFasthttpRequest(fastReq, false)
//
//	if err != nil {
//		return false
//	}
//	if bytes.Contains(resp.Body, existsStr) {
//		return true
//	}
//	return false
//}
