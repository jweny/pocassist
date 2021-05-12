package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"pocassist/basic"
	"pocassist/rule"
	"pocassist/utils"
	"sync"
)

var subCommandCli = cli.Command {
	Name:     "cli",
	Aliases:  []string{"c"},
	Usage:    "cli",
	Category: "cli",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "loadPoc",
			Aliases: []string{"lp"},
			Destination: &loadPoc,
			Value: "",
			Usage:
			"type of load plugins:\n" +
				"		case single: load one plugin\n" +
				"		case multi: load multi plugins\n" +
				"		case all: load all plugins(enable + disable)\n" +
				"		case affects: load one type `affects`\n" +
				"		default: load all enable plugins\n"},
		&cli.StringFlag{
			Name: "condition",
			Aliases: []string{"c"},
			Destination: &condition,
			Value: "",
			Usage:
			"the condition when set loadPoc:\n" +
				"		case single: condition is poc_id of plugins, example: `poc-db-001`\n" +
				"		case multi:, condition is multi poc_id of plugins, example: `poc-db-001,poc-db-002`\n" +
				"		case all:, condition is not use\n" +
				"		case affects, condition is name of `affects`, example: `server`\n" +
				"		case default, ``\n"},
		&cli.StringFlag{
			Name: "url",
			Aliases: []string{"u"},
			Destination: &url,
			Value: "",
			Usage: "only single url to scan"},
		&cli.StringFlag{
			Name: "urlFile",
			Aliases: []string{"uf"},
			Destination: &urlFile,
			Value: "",
			Usage: "load urls from file"},
		&cli.StringFlag{
			Name: "urlRaw",
			Aliases: []string{"ur"},
			Destination: &rawFile,
			Value: "",
			Usage: "load urls from request raw `File`"},
	},
	Action: RunCli,
}

func RunCli(c *cli.Context) error{
	err := InitAll()
	if err != nil {
		basic.GlobalLogger.Error("[init err ]", err)
	}
	// 加载poc
	vuls, err := rule.LoadPlugins(loadPoc, condition)
	if err != nil {
		basic.GlobalLogger.Error("[plugins load error ]" , err)
		return err
	}
	basic.GlobalLogger.Debug("[plugins load success]")

	// 并发限制
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(basic.GlobalConfig.PluginsConfig.Parallel, func(i interface{}) {
		rule.RunPoc(i)
		rule.ScanItemPut(i)
		wg.Done()
	})
	wg.Wait()
	defer p.Release()

	// 速率限制
	rule.InitRate()

	switch {
	case url != "":
		oreq, err := utils.GenOriginalReq(url)
		if err != nil {
			basic.GlobalLogger.Error("[original request gen err ]" , err)
			return err
		}
		basic.GlobalLogger.Info("[start check url ]" ,url)
		rule.RunPlugins(oreq, vuls)

	case urlFile != "":
		targets := utils.ReadingLines(urlFile)
		for _, url := range targets {
			oreq, err := utils.GenOriginalReq(url)
			if err != nil {
				basic.GlobalLogger.Error("[original request gen err ]" , err)
				return err
			}
			basic.GlobalLogger.Info("[start check url ]" ,url)
			rule.RunPlugins(oreq, vuls)
		}
	case rawFile != "":
		raw, err := ioutil.ReadFile(rawFile)
		if err != nil {
			basic.GlobalLogger.Error("[load url from file err ]" , err)
			return err
		}
		oreq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
		if err != nil {
			basic.GlobalLogger.Error("[original request gen err ]" , err)
			return err
		}
		if !oreq.URL.IsAbs() {
			scheme := "http"
			oreq.URL.Scheme = scheme
			oreq.URL.Host = oreq.Host
		}
		basic.GlobalLogger.Info("[start check url ]" ,oreq.URL.String())
		rule.RunPlugins(oreq, vuls)

	default:
		fmt.Println("Use -h for help")
	}
	return nil
}

