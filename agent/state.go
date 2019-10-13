package agent

import (
	"github.com/hpcloud/tail"
)

// Element: [offset, lastRead]
type logStatus map[string][2]int64

var lsCh = make(chan logStatus)
var tailCh = make(chan map[string]*tail.Tail)
var lsCtt = logStatus{}
var delCh = make(chan string)

func updateState() {
	for {
		select {
		case ls := <-lsCh:
			for key, val := range ls {
				lsCtt[key] = val
			}
		case tailData := <-tailCh:
			if tailData != nil {
				for k, v := range tailData {
					if tailIns == nil {
						tailIns = map[string]*tail.Tail{}
					}
					tailIns[k] = v
				}
			}
		case k := <-delCh:
			delete(lsCtt, k)
			delete(tailIns, k)
		}
	}
}
