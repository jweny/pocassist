package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Headers struct {
	UserAgent string `mapstructure:"user_agent"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Database string `mapstructure:"database"`
	Timeout  string `mapstructure:"timeout"`
}

type DbConfig struct {
	Sqlite string `mapstructure:"sqlite"`
	Mysql Mysql `mapstructure:"mysql"`
}

type PluginsConfig struct {
	Parallel int `mapstructure:"parallel"`
}

type Reverse struct {
	ApiKey string 	`mapstructure:"api_key"`
	Domain  string  `mapstructure:"domain"`
}

type Config struct {
	HttpConfig    HttpConfig    `mapstructure:"httpConfig"`
	DbConfig      DbConfig      `mapstructure:"dbConfig"`
	PluginsConfig PluginsConfig `mapstructure:"pluginsConfig"`
	Reverse       Reverse       `mapstructure:"reverse"`
	ServerConfig  ServerConfig	`mapstructure:"serverConfig"`
}

type ServerConfig struct {
	JwtSecret	string	`mapstructure:"jwt_secret"`
	RunMode		string	`mapstructure:"run_mode"`
	LogName		string	`mapstructure:"log_name"`

}

type HttpConfig struct {
	Headers     Headers `mapstructure:"headers"`
	Proxy       string  `mapstructure:"proxy"`
	HttpTimeout int     `mapstructure:"http_timeout"`
	DailTimeout int     `mapstructure:"dail_timeout"`
	UdpTimeout  int     `mapstructure:"udp_timeout"`
	MaxQps      int     `mapstructure:"max_qps"`
	MaxRedirect int     `mapstructure:"max_redirect"`
}

var GlobalConfig *Config

// 加载配置
func Setup() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("conf.Setup, fail to get current path: %v", err)
	}
	// 配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(dir, "config.yaml")

	// 检测配置文件是否存在
	_ , err = os.Lstat(configFile)
	if err != nil {
		//	没有，生成默认yaml
		WriteYamlConfig(configFile)
	}
	ReadYamlConfig(configFile)

}

func ReadYamlConfig(configFile string) {
	// 加载config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("conf.Setup, fail to read 'config.yaml': %v", err)
	}
	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.Fatalf("conf.Setup, fail to parse 'config.yaml', check format: %v", err)
	}
	err = verifiyConfig()
	if err != nil {
		log.Fatalf("conf.Setup, fail to verify 'config.yaml', check format: %v", err)
	}
}

func WriteYamlConfig(configFile string) {
	// 生成默认config
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultYamlByte))
	if err != nil {
		log.Fatalf("conf.Setup, fail to read default config bytes: %v", err)
	}
	// 写文件
	err = viper.SafeWriteConfigAs(configFile)
	if err != nil {
		log.Fatalf("conf.Setup, fail to write 'config.yaml': %v", err)
	}
}