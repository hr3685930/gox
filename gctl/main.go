package main

import (
	"github.com/hr3685930/pkg/gctl/project"
	"github.com/hr3685930/pkg/gctl/repo"
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

	repoFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "repo类型, 支持simple,db",
		},
		cli.StringFlag{
			Name:  "dir",
			Usage: "repo生成的路径, 没有则创建, repo名称根据最后一层级来命名",
		},
		cli.StringFlag{
			Name:  "model",
			Usage: "model名称, 需要放在models目录下, type为db时该字段生效",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "new",
			Usage:  "创建项目",
			Flags:  flags,
			Action: project.Create,
		},
		{
			Name:   "repo",
			Usage:  "创建repo",
			Flags:  repoFlags,
			Action: repo.Create,
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	_ = app.Run(os.Args)
}
