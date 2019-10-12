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
		for key, v := range allPath {
			if !v.alive {
				if time.Since(allPath[key].lastRead) > time.Second*time.Duration(aliveTime) {
					tailErr := v.tail.Stop()
					delete(allPath, key)
					util.ErrHandler(tailErr)
					// delete map data.
				}
			}
		}
	}

}
