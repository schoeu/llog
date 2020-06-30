package agent

import (
	"runtime"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/schoeu/llog/config"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

type storeState map[string]int64

var sm = cmap.New()

func StartAction(c *cli.Context) {
	defer util.Recover()
	launch(c.Args().First())
	select {}
}

func launch(args string) {
	var once sync.Once
	configFile := util.GetAbsPath(util.GetCwd(), args)
	err := config.InitCfg(configFile)
	conf := config.GetConfig()
	util.ErrHandler(err)

	taskInit(conf)

	// debug
	//debugInfo()

	inputs := conf.Input
	for _, v := range inputs {
		types := v.Type
		if types == "" {
			types = "log"
		}

		// type: log
		if types == "log" {
			logInput(v)
			once.Do(func() {
				// start watch file
				go watch(addWatchFile())

				if conf.SnapShot.Enable {
					// take snapshot for file status
					snd := conf.SnapShot.SnapShotDuring
					if snd == 0 {
						snd = snapShotDefault
					}
					takeSnap(snd)

					// recover file state
					recoverState()
				}
			})
		} else if types == "stdin" {
			// type: stdin
			stdInput()
		}
	}
}

func taskInit(conf *config.Config) {
	if conf.MaxProcs != 0 {
		runtime.GOMAXPROCS(conf.MaxProcs)
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
	info := conf.SysInfo
	if info {
		sysInfo(conf.SysInfoDuring)
	}
}
