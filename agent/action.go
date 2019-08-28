package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hpcloud/tail"
	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/nma/util"
	"github.com/urfave/cli"
)

func StartAction(c *cli.Context) error {
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	conf, err := util.GetConfig(configFile)
	logFile := conf.LogDir
	if logFile == "" {
		logFileDir, _ := os.UserHomeDir()
		logFile = path.Join(logFileDir, util.LogFileName)
	}
	t, err := tail.TailFile(logFile, tail.Config{
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
		Follow: true,
	})
	for line := range t.Lines {
		var psInfo gopsinfo.PsInfo
		if !conf.NoSysInfo {
			psInfo = gopsinfo.GetPsInfo(500)
		}
		var nodeInfo interface{}
		err = json.Unmarshal([]byte(line.Text), &nodeInfo)
		combineRs := util.CombineData(nodeInfo, psInfo, conf.NoSysInfo)
		PushData(combineRs)
	}
	return err
}

func StopAction(c *cli.Context) error {
	fmt.Println("stopAction")
	return nil
}

func RemoveAction(c *cli.Context) error {
	fmt.Println("removeAction")
	// TODO
	return nil
}

func StatusAction(c *cli.Context) error {
	fmt.Println("statusAction")
	// TODO
	return nil
}
