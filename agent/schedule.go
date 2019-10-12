package agent

import (
	"time"
)

func schedule(during int, fn func()) {
	ticker := time.NewTicker(time.Duration(during) * time.Second)
	for {
		<-ticker.C
		fn()
	}
}
