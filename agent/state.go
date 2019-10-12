package agent

import (
	"github.com/hpcloud/tail"
	"log"
	"time"
)

type allLogState struct {
	offset   int64
	lastRead time.Time
	alive    bool
	tail     *tail.Tail
}

var lsCh = make(chan []allLogState)

func updateState() {
	for data := range lsCh {
		log.Printf("You say: %s", data)
	}
}
