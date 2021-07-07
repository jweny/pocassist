package util

import (
	"crypto/tls"
	"github.com/jweny/pocassist/pkg/conf"
	"net"
	"time"
)
// TcpSend 指定目标发送tcp报文，返回结果（仅适用于一次交互即可判断漏洞的场景）
func TcpSend(targetAddr string, data []byte) ([]byte, error) {
	tcpTimeout := time.Duration(conf.GlobalConfig.HttpConfig.DailTimeout) * time.Second
	conn, err := net.DialTimeout("tcp", targetAddr, tcpTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(tcpTimeout))

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 20480)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}


// TcpSend 指定目标发送tls嵌套的tcp报文，返回结果（仅适用于一次交互即可判断漏洞的场景）
func TcpTlsSend(targetAddr string, data []byte) ([]byte, error) {
	tcpTimeout := time.Duration(conf.GlobalConfig.HttpConfig.HttpTimeout) * time.Second
	conn, err := net.DialTimeout("tcp", targetAddr, tcpTimeout)
	// add tls
	conn = net.Conn(tls.Client(conn, &tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 20480)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}