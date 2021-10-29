package scripts

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jweny/pocassist/pkg/util"
)

// RedisWeakPass redis 弱密码
func RedisWeakPass(args *ScriptScanArgs) (*util.ScanResult, error) {

	users := []string{"root", "redis", "admin", "default"}
	passes := []string{"123456", "root", "admin", "default", "redis", "root123"}

	version, err := judgeVersion(args.Host, 6379)
	if err != nil {
		return &util.InVulnerableResult, err
	}

	for _, user := range users {
		for _, pass := range passes {
			if version == "<6" {
				user = ""
			}
			flag, _ := conn(args.Host, 6379, user, "", 3)
			if flag {
				return util.VulnerableTcpOrUdpResult(args.Host+":6379", "Redis weak password.",
					[]string{string(user + ":" + pass)},
					[]string{},
				), nil
			}
		}

	}
	return &util.InVulnerableResult, nil
}

func init() {
	ScriptRegister("poc-go-redis-weakpass", RedisWeakPass)
}

func judgeVersion(host string, port int64) (string, error) {

	realhost := fmt.Sprintf("%s:%v", host, port)
	conn, err := net.DialTimeout("tcp", realhost, time.Duration(3)*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", "c4a6c07a8a2d7c804a5776d9d039428a")))
	if err != nil {
		return "", err
	}
	reply, err := readreply(conn)
	if err != nil {
		return "", err
	}
	if strings.Contains(reply, "invalid username-password pair or user is disabled") {

		return ">=6", nil
	}

	return "<6", nil
}

func readreply(conn net.Conn) (result string, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			break
		}
		result += string(buf[0:count])
		if count < 4096 {
			break
		}
	}
	return result, err
}

func conn(host string, port int64, user string, pass string, timeout int64) (bool, error) {
	realhost := fmt.Sprintf("%s:%v", host, port)
	conn, err := net.DialTimeout("tcp", realhost, time.Duration(timeout)*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s %s\r\n", user, pass)))
	if err != nil {
		return false, err
	}
	reply, err := readreply(conn)
	if err != nil {
		return false, err
	}
	if strings.Contains(reply, "+OK") {
		return true, nil
	}
	return false, nil
}
