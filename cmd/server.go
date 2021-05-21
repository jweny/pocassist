package cmd

import (
	"github.com/jweny/pocassist/api/routers"
	"github.com/urfave/cli/v2"
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
			Usage: "web server `PORT`",
		},
	},
	Action: RunServer,
}

func RunServer(c *cli.Context) error {
	InitAll()
	port := c.String("port")
	routers.InitRouter(port)
	return nil
}



