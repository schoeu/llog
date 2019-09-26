package main

import (
	"fmt"
	"github.com/schoeu/nma/agent"
	"github.com/schoeu/nma/util"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = util.Version
	app.Name = util.AppName
	app.Usage = util.AppUsage
	app.Action = agent.StartAction
	app.Commands = []cli.Command{
		{
			Name:   "stop",
			Usage:  "stop nma.",
			Action: agent.StopAction,
		},
	}
	err := app.Run(os.Args)
	util.ErrHandler(err)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}
