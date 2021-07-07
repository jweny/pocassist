package scripts

import (
	"bytes"
	"github.com/jweny/pocassist/pkg/util"
)

// RedisUnauthority redis 未授权 Poc
func RedisUnauthority(args *ScriptScanArgs) (*util.ScanResult, error) {
	addr := args.Host + ":6379"
	payload := []byte("*1\r\n$4\r\ninfo\r\n")
	resp, err := util.TcpSend(addr, payload)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(resp, []byte("redis_version")) {
		return util.VulnerableTcpOrUdpResult(addr, "",
			[]string{string(payload)},
			[]string{string(resp)},
		), nil
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-redis-unauth", RedisUnauthority)
}
