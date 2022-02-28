package main

import (
	"github.com/hr3685930/pkg/ctl/project"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func main() {
	app := cli.NewApp()
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "项目名称,和go mod同名",
		},
		cli.StringFlag{
			Name:  "err",
			Usage: "错误上报,支持sentry",
		},
		cli.StringFlag{
			Name:  "trace",
			Usage: "链路,支持jaeger",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "new",
			Usage:  "创建项目",
			Flags:  flags,
			Action: project.Create,
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	_ = app.Run(os.Args)
}
