package cmd

import (
	"github.com/jweny/pocassist/api/routers"
	conf2 "github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/urfave/cli/v2"
	"os"
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

func InitAll() {
	// config 必须最先加载
	conf2.Setup()
	logging.Setup(debug)
	db.Setup(dbname)
	routers.Setup()
	util.Setup()
	util.Setup()
	rule.Setup()
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
		logging.GlobalLogger.Error("[app run err ]", err)
		return
	}
}
