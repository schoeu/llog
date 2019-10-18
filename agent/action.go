package agent

import (
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
	"path/filepath"
)

func StartAction(c *cli.Context) {
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
