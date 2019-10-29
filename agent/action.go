package agent

import (
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

	//go updateState()

	for _, v := range inputs {
		// collect log.
		fileGlob(&v)
		// close file handle schedule.
		go closeFileHandle(&v)
	}
	addWatch()
	watch()

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

	ch := make(chan int)
	<-ch
}
