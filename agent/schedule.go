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
		timeoutDel <- aliveTime
	}
}
