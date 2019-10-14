package agent

import (
	"time"

	"github.com/schoeu/llog/util"
)

func scanFiles(fn func([]string), allLogs []string) {
	conf := util.GetConfig()
	sf := conf.ScanFrequency
	if sf < 1 {
		sf = 600
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
					err := v.Stop()
					util.ErrHandler(err)
					// delete map data.
					delCh <- key
				}
			}
		}
	}
}
