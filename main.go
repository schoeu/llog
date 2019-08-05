package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/pslog_agent/agent"
	"github.com/schoeu/pslog_agent/util"
	"github.com/urfave/cli"
)

type Config struct {
	AppId    string
	Secret   string
	LogDir   string
	Interval int
}

func main() {
	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "pslog_agent"
	app.Usage = "Agent for ps log"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "configuration file path.",
		},
	}

	app.Action = defaultAction
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start app on agent.",
			Action: startAction,
		},
		{
			Name:   "stop",
			Usage:  "stop app on agent.",
			Action: stopAction,
		},
		{
			Name:   "status",
			Usage:  "show app status.",
			Action: statusAction,
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "remove app.",
			Action:  removeAction,
		},
	}

	err := app.Run(os.Args)
	util.ErrHandler(err)
}

func getConfig(p string) (Config, error) {
	p = util.GetAbsPath(util.GetHomeDir(), p)

	c := Config{}
	data, err := ioutil.ReadFile(p)
	err = json.Unmarshal(data, &c)

	return c, err
}

func removeAction(c *cli.Context) error {
	fmt.Println("removeAction")
	return nil
}

func statusAction(c *cli.Context) error {
	fmt.Println("statusAction")
	return nil
}

func startAction(c *cli.Context) error {
	fmt.Println("startAction")
	return nil
}

func stopAction(c *cli.Context) error {
	fmt.Println("stopAction")
	return nil
}

func defaultAction(c *cli.Context) error {
	configFile := util.GetAbsPath("", c.Args().First())
	ext := path.Ext(configFile)
	if ext == ".json" {
		conf, err := getConfig(configFile)
		util.ErrHandler(err)
		psInfoTimer(conf)

	} else {
		fmt.Println("Invited json file.")
	}

	return nil
}

func psInfoTimer(conf Config) {
	d := time.Duration(time.Millisecond * time.Duration(conf.Interval))
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		psInfo := gopsinfo.GetPsInfo(conf.Interval)
		agent.PushData(&psInfo,  conf.AppId, conf.Secret)
	}
}
