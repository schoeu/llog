package agent

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
	"runtime"
)

type storeState map[string][2]int64

var sm = cmap.New()

func StartAction(c *cli.Context) {
	defer util.Recover()

	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	err := util.InitCfg(configFile)
	conf := util.GetConfig()
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

		// TODO
		reScanTask(&v)
	}

	for _, v := range sm.Keys() {
		li, _ := getLogInfoIns(v)
		fmt.Println("sm.Keys()111", v, li.Sc.ScanFrequency)
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
	addWatch()

	// watch file change
	watch()

	if conf.SnapShot.Enable {
		// take snapshot for file status
		takeSnap()

		// recover file state
		recoverState()
	}

	// debug
	debugInfo()
}

func reScan() {
	inputs := util.GetConfig().Input
	for _, v := range inputs {
		// collect log.
		fileGlob(v)
	}
}
