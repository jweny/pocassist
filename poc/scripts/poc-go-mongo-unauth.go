package scripts

import (
	"net"
	"strings"
	"time"

	"github.com/jweny/pocassist/pkg/util"
)

// MongoDBUnauthority MongoDB 未授权
func MongoDBUnauthority(args *ScriptScanArgs) (*util.ScanResult, error) {
	addr := args.Host + ":27017"
	senddata := []byte{58, 0, 0, 0, 167, 65, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 255, 255, 255, 255, 19, 0, 0, 0, 16, 105, 115, 109, 97, 115, 116, 101, 114, 0, 1, 0, 0, 0, 0}
	getlogdata := []byte{72, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 1, 0, 0, 0, 33, 0, 0, 0, 2, 103, 101, 116, 76, 111, 103, 0, 16, 0, 0, 0, 115, 116, 97, 114, 116, 117, 112, 87, 97, 114, 110, 105, 110, 103, 115, 0, 0}
	payload := append(senddata, getlogdata...)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(senddata)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1024)
	count, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	text := string(buf[0:count])
	if strings.Contains(text, "ismaster") {
		_, err = conn.Write(getlogdata)
		if err != nil {
			return nil, err
		}
		count, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		text := string(buf[0:count])
		if strings.Contains(text, "totalLinesWritten") {
			return util.VulnerableTcpOrUdpResult(addr, "",
				[]string{string(payload)},
				[]string{text},
			), nil
		}
	}

	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-mongo-unauth", MongoDBUnauthority)
}
