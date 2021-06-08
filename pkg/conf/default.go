package conf

import (
	"encoding/json"
	"errors"
)

// todo 改回 release
var defaultYamlByte = []byte(`
# webserver配置
serverConfig:
  # 配置jwt秘钥
  jwt_secret: "pocassist"
  # gin的运行模式 "release" 或者 "debug"
  run_mode: "debug"
  # 运行日志的文件名，日志将保存在二进制所在目录
  log_name : "debug.log"

# HTTP配置
httpConfig:
  # 扫描时使用的代理：格式为 IP:PORT，example: 如 burpsuite，可填写 127.0.0.1:8080
  proxy: ""
  # 读取 http 响应超时时间，不建议设置太小，否则可能影响到盲注的判断
  http_timeout: 10
  # 建立 tcp 连接的超时时间
  dail_timeout: 5
  # udp 超时时间
  udp_timeout: 5
  # 每秒最大请求数
  max_qps: 100
  # 单个请求最大允许的跳转次数
  max_redirect: 5
  headers:
    # 默认 UA
    user_agent: "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0"

# 数据库配置
dbConfig:
  # sqlite配置：sqlite数据库文件的路径
  sqlite : "pocassist.db"
  # mysql配置
  mysql:
    host: "127.0.0.1"
    password: ""
    port: "3306"
    user: "root"
    database: "pocassist"
    # 数据库连接超时时间
    timeout: "3s"

# 插件配置
pluginsConfig:
  # 并发量:同时运行的插件数量
  parallel: 8

# 反连平台配置: 目前使用 ceye.io
reverse:
  api_key: ""
  domain: ""
`)

var runMode = []string{"debug","release"}

func ArrayToString (array []string) string {
	str, _ := json.Marshal(array)
	return string(str)
}

func StrInArray (str string, array []string) error {
	for _, element := range array{
		if str == element{
			return nil
		}
	}
	return errors.New(str + "must in" + ArrayToString(array))
}

func verifiyConfig() error {
	var err error
	err = StrInArray(GlobalConfig.ServerConfig.RunMode, runMode)
	if err != nil {
		return err
	}
	return nil
}
