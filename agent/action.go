package agent

import (
	"fmt"
	"path/filepath"

	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

var (
	gConf   util.Config
	allLogs []string
)

func StartAction(c *cli.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	conf, err := util.GetConfig(configFile)
	gConf = conf
	util.ErrHandler(err)

	logFiles := conf.LogDir

	var ch = make(chan int)

	if len(logFiles) == 0 {
		logFileDir := util.GetTempDir()
		logFiles = append(logFiles, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	allLogs = logFiles
	// 监控日志收集
	fileGlob()

	util.ErrHandler(err)

	sf := gConf.ScanFrequency
	if sf == 0 {
		sf = 10
	}

	// log file scan schedule.
	scan(sf, fileGlob)

	<-ch
}

//func StopAction(c *cli.Context) {
//	fmt.Println("stop")
//}
