package agent

import (
	"github.com/schoeu/llog/config"
	"runtime"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

type storeState map[string][2]int64

var sm = cmap.New()

func StartAction(c *cli.Context) {
	defer util.Recover()

	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	err := config.InitCfg(configFile)
	conf := config.GetConfig()
	util.ErrHandler(err)

	if conf.MaxProcs != 0 {
		runtime.GOMAXPROCS(conf.MaxProcs)
	}

	inputs := conf.Input
	for _, v := range inputs {
		// collect log
		fileGlob(v)

		// close file handle schedule
		closeFileHandle(v)

		// watch new log file schedule
		reScanTask(v.ScanFrequency)
	}

	// set app name
	appName := conf.Name
	if appName == "" {
		appName = util.AppName
	}
	name = appName

	// set api server info
	output := conf.Output
	apiEnable := output.ApiServer.Enable
	if apiEnable && output.ApiServer.Url != "" {
		apiServer = output.ApiServer.Url
	}

	// init es
	es := conf.Output.Elasticsearch
	if es.Enable && len(es.Host) > 0 {
		esInit()
	}

	// system info process
	sysInfo()

	// start watch file
	fsWatcher := addWatchFile()

	// watch file change
	watch(fsWatcher)

	if conf.SnapShot.Enable {
		// take snapshot for file status
		takeSnap()

		// recover file state
		recoverState()
	}

	// debug
	//debugInfo()
}
