package main

import (
	"log"
	"os"
	"path"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Version = "1.0.0"
	app.Name = "pslog_agent"
	app.Usage = "Agent for ps log"
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Usage:   "start app on agent.",
			Action:  startAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startAction(c *cli.Context) error {
	//fmt.Println("added task: ", c.Args())


	return nil
}

func getLog() {
	defaultLogName := ".psinfo.log"
	homeDir := os.UserHomeDir()
	defaultLogFile := path.Join(homeDir, defaultLogName)


}