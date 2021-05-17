package basic

import (
	"github.com/spf13/viper"
	"path"
)

// Headers
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

// DbConfig
type DbConfig struct {
	Sqlite string `mapstructure:"sqlite"`
	Mysql Mysql `mapstructure:"mysql"`
}

// PluginsConfig
type PluginsConfig struct {
	Parallel int `mapstructure:"parallel"`
}

// Reverse
type Reverse struct {
	ApiKey string 	`mapstructure:"api_key"`
	Domain  string  `mapstructure:"domain"`
}

// Yaml2Go
type Config struct {
	HttpConfig    HttpConfig    `mapstructure:"httpConfig"`
	DbConfig      DbConfig      `mapstructure:"dbConfig"`
	PluginsConfig PluginsConfig `mapstructure:"pluginsConfig"`
	Reverse       Reverse       `mapstructure:"reverse"`
}

// HttpConfig
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

// 加载配置 & log
func InitConfig(workDir string) error {
	// 加载config
	configFile := path.Join(workDir, "config.yaml")
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		GlobalLogger.Error("[init fail config.yaml read err ]" + err.Error())
		return err
	}
	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		GlobalLogger.Error("[init fail config.yaml load err ]" + err.Error())
		return err
	}
	return nil
}


