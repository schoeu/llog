package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/urfave/cli"
	"github.com/schoeu/pslog_agent/util"
	"github.com/takama/daemon"
)

type config struct {
	Appid  string
	Secret string
	Logdir string
	Interval int
}

var (
	service daemon.Daemon
)

func main() {
	app := cli.NewApp()

	service, _ = daemon.New("pslogAgent", "description")

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
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "status fro agent.",
			Action:  listAction,
		},
		{
			Name:    "stop",
			Usage:   "stop app on agent.",
			Action:  stopAction,
		},
		{
			Name:    "status",
			Usage:   "show app status.",
			Action:  statusAction,
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

func getConfig(p string) (config, error) {
	p = util.GetAbsPath(util.GetHomeDir(), p)

	c := config{}
	data, err := ioutil.ReadFile(p)
	err = json.Unmarshal(data, &c)

	return c, err
}

func removeAction(c *cli.Context) (string, error) {
	return service.Remove()
}

func statusAction(c *cli.Context) (string, error) {
	return service.Status()
}

func startAction(c *cli.Context)  (string, error) {
	return service.Start()
}

func listAction(c *cli.Context) error {
	fmt.Println("listAction", c)
	return nil
}

func stopAction(c *cli.Context)  (string, error) {
	return service.Stop()
}

func defaultAction(c *cli.Context) error {
	configFile := c.Args().First()
	configFile = util.GetAbsPath("", configFile)
	ext := path.Ext(configFile)
	if ext == ".json" {
		conf, err := getConfig(configFile)
		util.ErrHandler(err)
		psInfoTimer(conf.Interval)

	} else {
		fmt.Println("Invited json file.")
	}

	status, err := service.Install()
	util.ErrHandler(err)
	fmt.Println(status)

	return nil
}

func psInfoTimer(interval int) {
	d := time.Duration(time.Millisecond * time.Duration(interval))
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		psInfo := gopsinfo.GetPsInfo(interval)
		fmt.Println(psInfo)
	}
}
