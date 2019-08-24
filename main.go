package main

import (
	"fmt"
	"os"

	"github.com/schoeu/nma/agent"
	"github.com/schoeu/nma/util"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = util.Version
	app.Name = util.AppName
	app.Usage = util.AppUsage

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "configuration file path.",
		},
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	app.Action = agent.StartAction
	app.Commands = []cli.Command{
		{
			Name:   "stop",
			Usage:  "stop app on agent.",
			Action: agent.StopAction,
		},
		{
			Name:   "status",
			Usage:  "show app status.",
			Action: agent.StatusAction,
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "remove app.",
			Action:  agent.RemoveAction,
		},
	}

	err := app.Run(os.Args)
	util.ErrHandler(err)
}
