package agent

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hpcloud/tail"
	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/nma/util"
	"github.com/urfave/cli"

)

func StartAction(c *cli.Context) error {
	configFile := util.GetAbsPath("", c.Args().First())
	conf, err := util.GetConfig(configFile)
	t, err := tail.TailFile(conf.LogDir, tail.Config{
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
		Follow: true,
	})
	for line := range t.Lines {
		psInfo := gopsinfo.GetPsInfo(100)
		var nodeInfo interface{}
		err = json.Unmarshal([]byte(line.Text), &nodeInfo)
		combineRs := util.CombineData(nodeInfo, psInfo)
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
