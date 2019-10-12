package agent

import (
	"fmt"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
	"path/filepath"
)

var (
	allLogs []string
)

func StartAction(c *cli.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	err := util.InitCfg(configFile)
	conf := util.GetConfig()
	util.ErrHandler(err)

	logFiles := conf.LogDir

	if len(logFiles) == 0 {
		logFileDir := util.GetTempDir()
		logFiles = append(logFiles, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	allLogs = logFiles

	// 监控日志收集
	fileGlob()

	util.ErrHandler(err)

	sf := conf.ScanFrequency
	if sf < 1 {
		sf = 10
	}
	// log file scan schedule.
	go schedule(sf, fileGlob)
	ci := conf.CloseInactive
	if ci < 1 {
		ci = 300
	}
	fmt.Println("aa")
	schedule(ci, closeFileHandle)
}

//func StopAction(c *cli.Context) {
//	fmt.Println("stop")
//}
