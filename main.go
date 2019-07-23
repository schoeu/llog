package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/urfave/cli"
)

type config struct {
	Appid  string
	Secret string
	Logdir string
}

var (
	//interval = 60000
	interval = 1000
)

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
		log.Fatal(err)
	}

}

func getLog() {
	defaultLogName := ".psinfo.log"
	homeDir, err := os.UserHomeDir()
	errHandler(err)
	defaultLogFile := path.Join(homeDir, defaultLogName)
	fmt.Println(defaultLogFile)

}

func getConfig(p string) (config, error) {
	if !path.IsAbs(p) {
		p = path.Join(getCwd(), p)
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
	fmt.Println("222", configFile)
	if ext == ".json" {
		conf, err := getConfig(configFile)
		errHandler(err)

		// get config data: Appid, Secret, Logdir
		fmt.Println("defaultAction", conf.Appid)
	}
	return nil
}

func timer(interval int) {
	d := time.Duration(time.Millisecond * time.Duration(interval))
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		getLog()
	}
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
