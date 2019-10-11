package agent

import (
	"time"
)

func scan(during int, scanFunc func()) {
	ticker := time.NewTicker(time.Duration(during) * time.Second)
	for {
		<-ticker.C
		scanFunc()
	}
}
