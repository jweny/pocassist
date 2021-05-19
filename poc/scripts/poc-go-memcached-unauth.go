package scripts

import (
	"bytes"
	"github.com/jweny/pocassist/pkg/util"
)

// MemcachedUnauthority memcached 未授权
func MemcachedUnauthority(args *ScriptScanArgs) (*util.ScanResult, error) {
	addr := args.Host + ":11211"
	payload := []byte("stats\n")
	resp, err := util.TcpSend(addr, payload)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(resp, []byte("STAT pid")) {
		return util.VulnerableTcpOrUdpResult(addr, "",
			[]string{string(payload)},
			[]string{string(resp)},
		),nil
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-memcached-unauth", MemcachedUnauthority)
}
