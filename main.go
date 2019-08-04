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
)

type config struct {
	Appid  string
	Secret string
	Logdir string
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
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "status fro agent.",
			Action:  listAction,
		},
		{
			Name:    "stop",
			Aliases: []string{"l"},
			Usage:   "stop app on agent.",
			Action:  stopAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func getConfig(p string) (config, error) {
	if !path.IsAbs(p) {
		p = path.Join(util.GetHomeDir(), p)
	}

	c := config{}
	data, err := ioutil.ReadFile(p)
	err = json.Unmarshal(data, &c)

	return c, err
}

func startAction(c *cli.Context) error {
	fmt.Println("config")

	return nil
}

func listAction(c *cli.Context) error {
	fmt.Println("listAction", c)
	return nil
}

func stopAction(c *cli.Context) error {
	fmt.Println("stopAction", c)
	return nil
}

func defaultAction(c *cli.Context) error {
	configFile := c.Args().First()
	if !path.IsAbs(configFile) {
		configFile = path.Join(util.GetCwd(), configFile)
	}
	ext := path.Ext(configFile)
	if ext == ".json" {
		conf, err := getConfig(configFile)
		util.ErrHandler(err)
		fmt.Println(conf.Appid)
		psInfoTimer(conf.Interval)

		// get config data: Appid, Secret, Logdir
	} else {
		fmt.Println("Invited json file.")
	}
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
