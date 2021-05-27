package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jweny/pocassist/api/routers"
	conf2 "github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var (
	url			string
	urlFile		string
	rawFile		string
	loadPoc		string
	condition	string
)

func init() {
	welcome := `
                               _     _
 _ __   ___   ___ __ _ ___ ___(_)___| |_
| '_ \ / _ \ / __/ _' / __/ __| / __| __|
| |_) | (_) | (_| (_| \__ \__ \ \__ \ |_
| .__/ \___/ \___\__,_|___/___/_|___/\__|
|_|
`
	fmt.Println(welcome)
}

func InitAll() {
	// config 必须最先加载
	conf2.Setup()
	logging.Setup()
	db.Setup()
	routers.Setup()
	util.Setup()
	rule.Setup()
}

// 使用viper 对配置热加载
func HotConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("conf.Setup, fail to get current path: %v", err)
	}
	// 配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(dir, "config.yaml")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	// watch 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		log.Println("Config file changed:", e.Name)
		InitAll()
	})
}

func RunApp() {
	app := cli.NewApp()
	app.Name = "pocassist"
	app.Usage = "New POC Framework Without Writing Code"
	app.Version = "0.3.0"

	// 子命令
	app.Commands = []*cli.Command{
		&subCommandCli,
		&subCommandServer,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("cli.RunApp err: %v", err)
		return
	}
}
