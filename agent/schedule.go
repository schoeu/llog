package agent

import (
	"time"

	"github.com/schoeu/llog/util"
)

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
