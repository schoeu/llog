package agent

import (
	"fmt"
	"github.com/hpcloud/tail"
)

//type allLogState struct {
//	offset   int64
//	lastRead time.Time
//	alive    bool
//	tail     *tail.Tail
//}

type allLogState map[string]interface{}

var lsCh = make(chan map[string]allLogState)
var tailCh = make(chan map[string]*tail.Tail)

func updateState() {
	for {
		select {
		case ls := <-lsCh:
			for data := range ls {
				fmt.Println("ls", data)
			}
		case tailData := <-tailCh:
			if tailData != nil {
				for k, v := range tailData {
					tailIns[k] = v
					fmt.Println("tailIns", tailIns)
				}
			}
		}
	}
}
