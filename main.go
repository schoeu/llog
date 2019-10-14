package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/schoeu/llog/agent"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)
import _ "net/http/pprof"

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
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
