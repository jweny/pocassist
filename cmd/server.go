package cmd

import (
	"github.com/urfave/cli/v2"
	"pocassist/api"
	"pocassist/basic"
)

var subCommandServer = cli.Command{
	Name:     "server",
	Aliases:  []string{"s"},
	Usage:    "server",
	Category: "server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			// 后端端口
			Name:  "port",
			Aliases: []string{"p"},
			Value: "1231",
			Usage: "api server port",
		},
	},
	Action: RunServer,
}

func RunServer(c *cli.Context) error {
	err := InitAll()
	if err != nil {
		basic.GlobalLogger.Error("[init err ]", err)
		return err
	}
	port := c.String("port")
	api.Route(port)
	return nil
}



