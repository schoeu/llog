package main

import (
	"fmt"
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

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	app.Action = agent.StartAction
	//app.Commands = []cli.Command{
	//	{
	//		Name:   "stop",
	//		Usage:  "stop llog.",
	//		Action: agent.StopAction,
	//	},
	//}
	err := app.Run(os.Args)
	util.ErrHandler(err)
}
