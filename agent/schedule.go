package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

const aliveTimeDefault = 300
const freqDefault = 600

func closeFileHandle(sc *util.SingleConfig) {
	defer util.Recover()

	aliveTime := sc.CloseInactive
	if aliveTime < 1 {
		aliveTime = aliveTimeDefault
	}
	ticker := time.NewTicker(time.Duration(aliveTime) * time.Second)
	for {
		<-ticker.C
		for _, v := range sm.Keys() {
			li, err := getLogInfoIns(v)
			util.ErrHandler(err)
			if li != nil && li.sc == sc && time.Since(time.Unix(li.status[1], 0)) > time.Second*time.Duration(aliveTime) {
				fmt.Println("stop watch: ", v)
				delInfo(v)
			}
		}
	}
}

func reScanTask(sc *util.SingleConfig) {
	defer util.Recover()

	freq := sc.ScanFrequency
	if freq < 1 {
		freq = freqDefault
	}
	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	for {
		<-ticker.C
		reScan()
	}
}

func sysInfo() {
	conf := util.GetConfig()
	info := conf.SysInfo

	if info {
		during := conf.SysInfoDuring
		var psInfo gopsinfo.PsInfo
		var d time.Duration
		if during < 1 {
			d = 1
		} else if during == 0 {
			d = 10
		}
		ticker := time.NewTicker(d * time.Second)
		go func() {
			defer util.Recover()
			for {
				<-ticker.C

				psInfo = gopsinfo.GetPsInfo(d)
				sysData, err := json.Marshal(psInfo)
				util.ErrHandler(err)
				doPush(&sysData, systemType)
			}
		}()
	}
}

func debugInfo() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C
			fmt.Println("\n\n\nsm count: ", sm.Count())
			for k, v := range sm.Items() {
				val := v.(logInfo)
				fmt.Println(k, "--->", val)
			}
		}
	}()
}
