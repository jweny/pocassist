package scripts

import (
	"time"

	"github.com/jweny/pocassist/pkg/util"

	"github.com/jlaffaye/ftp"
)

// FtpUnauthority Ftp 未授权
func FtpUnauthority(args *ScriptScanArgs) (*util.ScanResult, error) {
	addr := args.Host + ":21"
	payload := "anonymous"
	con, err := ftp.DialTimeout(addr, 5*time.Second)

	if err == nil {
		err = con.Login("anonymous", "")
		if err == nil {
			defer con.Logout()
			return util.VulnerableTcpOrUdpResult(addr, "",
				[]string{string(payload)},
				[]string{},
			), nil
		}
	} else {
		return nil, err
	}

	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-ftp-unauth", FtpUnauthority)
}
