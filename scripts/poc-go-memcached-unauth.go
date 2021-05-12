package scripts

import (
	"bytes"
	"pocassist/utils"
)

// MemcachedUnauthority memcached 未授权
func MemcachedUnauthority(args *ScriptScanArgs) (*utils.ScanResult, error) {
	addr := args.Host + ":11211"
	payload := []byte("stats\n")
	resp, err := utils.TcpSend(addr, payload)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(resp, []byte("STAT pid")) {
		return utils.VulnerableTcpOrUdpResult(addr, "",
			[]string{string(payload)},
			[]string{string(resp)},
		),nil
	}
	return &utils.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-memcached-unauth", MemcachedUnauthority)
}
