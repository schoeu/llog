package main

import (
	"os"

	"github.com/schoeu/llog/agent"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = util.Version
	app.Name = util.AppName
	app.Usage = util.AppUsage

	defer util.Recover()

	app.Action = agent.StartAction
	err := app.Run(os.Args)
	util.ErrHandler(err)
}
