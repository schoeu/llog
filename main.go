package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/urfave/cli"
	"github.com/schoeu/gopsinfo"
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

func getPsInfo(during int) {
	psInfo := gopsinfo.GetPsInfo(during)
	fmt.Println(psInfo)
}

//func getTailData(val string) {
//	t, err := tail.TailFile(val, tail.Config{Follow: true})
//	errHandler(err)
//	for line := range t.Lines {
//		fmt.Println(val, "--->", line.Text)
//	}
//}
//
//func getLogPath(p string) string {
//	if p == "" {
//		defaultLogName := ".psinfo.log"
//		p = path.Join(getHomeDir(), defaultLogName)
//	}
//	return p
//}

func getConfig(p string) (config, error) {

	if !path.IsAbs(p) {
		p = path.Join(getHomeDir(), p)
	}

	c := config{}
	data, err := ioutil.ReadFile(p)
	err = json.Unmarshal(data, &c)

	return c, err
}

func errHandler(err error) {
	if err != nil {
		fmt.Println(err)
	}
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
	ext := path.Ext(configFile)
	if ext == ".json" {
		conf, err := getConfig(configFile)
		errHandler(err)
		fmt.Println(conf.Appid)
		//conf.Logdir
		//conf.Secret

		// get config data: Appid, Secret, Logdir
		logPath := getLogPath(conf.Logdir)
		watcher(logPath)

	} else {
		fmt.Println("Invited json file.")
	}
	return nil
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func getHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return dir
}
