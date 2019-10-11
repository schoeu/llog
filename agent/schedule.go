package agent

import (
	"time"
)

func scan(during int, fn func()) {
	ticker := time.NewTicker(time.Duration(during) * time.Second)
	for {
		<-ticker.C
		fn()
	}
}
