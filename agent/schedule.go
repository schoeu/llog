package agent

import (
	"github.com/schoeu/llog/util"
	"time"
)

func scanFiles(fn func([]string), allLogs []string) {
	conf := util.GetConfig()
	sf := conf.ScanFrequency
	if sf < 1 {
		sf = 10
	}
	ticker := time.NewTicker(time.Duration(sf) * time.Second)
	for {
		<-ticker.C
		fn(allLogs)
	}
}

func closeFileHandle() {
	conf := util.GetConfig()
	aliveTime := conf.CloseInactive
	if aliveTime < 1 {
		aliveTime = 300
	}
	ticker := time.NewTicker(time.Duration(aliveTime) * time.Second)
	for {
		<-ticker.C
		for key, v := range tailIns {
			if v != nil {

				if time.Since(time.Unix(lsCtt[key][1], 0)) > time.Second*time.Duration(aliveTime) {
					tailErr := v.Stop()
					delCh <- key
					util.ErrHandler(tailErr)
					// delete map data.
				}
			}
		}
	}
}
