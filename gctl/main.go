package main

import (
	"github.com/hr3685930/pkg/gctl/project"
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
		cli.StringFlag{
			Name:  "db",
			Usage: "数据库,支持mysql,postgre,clickhouse",
		},
		cli.StringFlag{
			Name:  "queue",
			Usage: "队列,支持kafka,rabbitmq",
		},
		cli.StringFlag{
			Name:  "cache",
			Usage: "缓存,支持redis",
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
