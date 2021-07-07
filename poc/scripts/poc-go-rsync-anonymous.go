package scripts

import (
	"bytes"
	"context"
	"os/exec"
	"github.com/jweny/pocassist/pkg/util"
	"strings"
	"time"
)

func runCommand(name string, params ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd := exec.CommandContext(ctx, name, params...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RsyncAnonymousAccess rsync 匿名访问
func RsyncAnonymousAccess(args *ScriptScanArgs) (*util.ScanResult, error) {
	stdout, stderr, err := runCommand("rsync", args.Host+"::")
	if err != nil {
		return nil, err
	}

	if len(stderr) == 0 {
		stdout = strings.Replace(strings.Replace(stdout, " ", "", -1), "\n", "", -1)
		if len(stdout) > 0 {
			path := strings.SplitN(stdout, "\t", 2)[0]
			command := args.Host + "::" + path
			stdout, stderr, err = runCommand("rsync", command)
			if err != nil {
				return nil, err
			}
			if len(stderr) == 0 {

				return util.VulnerableTcpOrUdpResult(command, "rsync " + command,
					nil,
					[]string{stdout},
				),nil
			}
		}
	}

	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-rsync-anonymous", RsyncAnonymousAccess)
}
