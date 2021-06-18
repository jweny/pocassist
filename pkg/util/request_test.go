package util

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"testing"
)


func TestDoFasthttpRequest(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://www.360.cn")
	requestBody := []byte(`{"request":"test"}`)
	req.SetBody(requestBody)

	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	url1 := string(req.RequestURI())
	url2 := req.URI().String()
	url3 := string(req.Host())
	url4 := string(req.URI().RequestURI())
	protocol := string(req.Header.Protocol())
	url6 := string(req.Header.Header())

	//absRequestURI := strings.HasPrefix(reqURI, "http://") || strings.HasPrefix(reqURI, "https://")
	fmt.Println(url1)
	fmt.Println(url2)
	fmt.Println(url3)
	fmt.Println(url4)
	fmt.Println(protocol)
	fmt.Println(url6)
	/*
	https://www.360.cn
	https://www.360.cn/
	www.360.cn
	/
	*/

}
