package agent

import (
	"path/filepath"
	"runtime"

	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

func StartAction(c *cli.Context) {
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	err := util.InitCfg(configFile)
	conf := util.GetConfig()
	util.ErrHandler(err)

	if conf.MaxProcs != 0 {
		runtime.GOMAXPROCS(conf.MaxProcs)
	}

	inputs := conf.Input

	for _, v := range inputs {
		logFiles := v.LogDir
		if len(logFiles) == 0 {
			logFileDir := util.GetTempDir()
			logFiles = append(logFiles, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
		}
		// collect log.
		fileGlob(logFiles)
	}

	util.ErrHandler(err)

	// init es
	if conf.Output.Elasticsearch.Enable && len(conf.Output.Elasticsearch.Host) > 0 {
		esInit()
	}

	// close file handle schedule.
	closeFileHandle()
}
