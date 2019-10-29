package agent

import (
	"encoding/json"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

func closeFileHandle(sc *util.SingleConfig) {
	aliveTime := sc.CloseInactive
	if aliveTime < 1 {
		aliveTime = 300
	}
	ticker := time.NewTicker(time.Duration(aliveTime) * time.Second)
	for {
		<-ticker.C
		for _, v := range sm.Keys() {
			ins, ok := sm.Get(v)
			if !ok {
				util.ErrHandler(syncMapError)
				return
			}
			li := ins.(logInfo)
			if time.Since(time.Unix(li.status[1], 0)) > time.Second*time.Duration(aliveTime) {
				delInfo(v)
			}
		}
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
