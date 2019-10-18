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

	// collect log.
	fileGlob(logFiles)

	util.ErrHandler(err)

	// init es
	if conf.Elasticsearch.Enable && len(conf.Elasticsearch.Host) > 0 {
		esInit()
	}

	// close file handle schedule.
	closeFileHandle()
}

//func StartAction1(c *cli.Context) {
//	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
//	err := util.InitCfg(configFile)
//	util.ErrHandler(err)
//	tailFile([]string{"/var/folders/lp/jd6nj9ws5r3br43_y7qw66zw0000gn/T/.nm_logs/nm_apps3/nm_app_nmt.log"})
//	var s = make(chan string)
//	<-s
//}
