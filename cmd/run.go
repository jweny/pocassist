package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"path/filepath"
	"pocassist/basic"
	"pocassist/database"
	"pocassist/rule"
	"pocassist/utils"
	"sort"
)

var (
	url			string
	urlFile		string
	rawFile		string
	loadPoc		string
	condition	string
	debug		bool
	dbname		string
)

func InitAll() error {
	// 获取当前工作目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("init fail get work dir err:" + err.Error())
		return err
	}
	// 初始化logger
	logFile := path.Join(dir, "debug.log")
	err = basic.InitLog(debug, logFile)
	if err != nil {
		return err
	}
	basic.GlobalLogger.Debug("[globalLogger init success]")
	// 加载配置文件
	err = basic.InitConfig(dir)
	if err != nil {
		basic.GlobalLogger.Debug("[config.yaml init fail]")
		return err
	}
	basic.GlobalLogger.Debug("[config.yaml init success]")
	// 建立数据库连接
	err = database.InitDB(dbname)
	if err != nil {
		return err
	}
	basic.GlobalLogger.Debug("[database init success]")
	// fasthttp client 初始化
	err = utils.InitFastHttpClient(basic.GlobalConfig.HttpConfig.Proxy)
	if err != nil {
		basic.GlobalLogger.Error("[fasthttp client init err ]", err)
		return err
	}
	basic.GlobalLogger.Debug("[fasthttp client init success]")
	// handle初始化
	rule.InitHandles()
	basic.GlobalLogger.Debug("[rule handles init success]")
	return nil
}

func RunApp() {
	app := cli.NewApp()
	app.Name = "PocAssist"
	app.Usage = "New POC Framework Without Writing Code"
	app.Version = "1.0.0"
	// 全局flag
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "debug",
			Aliases: []string{"d"},
			Destination: &debug,
			Value: false,
			Usage: "enable debug log"},
		&cli.StringFlag{
			Name: "database",
			Aliases: []string{"db"},
			Destination: &dbname,
			Value: "sqlite",
			Usage: "kind of database, default: sqlite"},
	}



	// 子命令
	app.Commands = []*cli.Command{
		&subCommandCli,
		&subCommandServer,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		basic.GlobalLogger.Error("[app run err ]", err)
		return
	}
}
