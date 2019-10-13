package agent

import (
	"fmt"

	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
	"path/filepath"
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
	go updateState()

	// collect log.
	fileGlob(logFiles)

	util.ErrHandler(err)

	// log file scan schedule.
	go scanFiles(fileGlob, logFiles)

	// close file handle schedule.
	closeFileHandle()
}

//func StopAction(c *cli.Context) {
//	fmt.Println("stop")
//}
