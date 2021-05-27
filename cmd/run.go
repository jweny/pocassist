package cmd

import (
	"github.com/jweny/pocassist/api/routers"
	conf2 "github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

var (
	url			string
	urlFile		string
	rawFile		string
	loadPoc		string
	condition	string
)

func InitAll() {
	// config 必须最先加载
	conf2.Setup()
	logging.Setup()
	db.Setup()
	routers.Setup()
	util.Setup()
	rule.Setup()
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
