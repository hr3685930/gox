package command

import (
	"flag"
	"github.com/urfave/cli"
	"os"
	"sort"
)

type Command struct {
	Commands []cli.Command
}

func NewCommand(commands []cli.Command) *Command {
	return &Command{Commands: commands}
}


func (c *Command) Init(){
	flag.Parse()
	if len(flag.Args()) == 0 {
		return
	}
	app := cli.NewApp()
	app.Commands = c.Commands
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	_ = app.Run(os.Args)
	os.Exit(0)
}