package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

const aliveTimeDefault = 300
const freqDefault = 600

func closeFileHandle(sc *util.SingleConfig) {
	aliveTime := sc.CloseInactive
	if aliveTime < 1 {
		aliveTime = aliveTimeDefault
	}
	ticker := time.NewTicker(time.Duration(aliveTime) * time.Second)

	go func() {
		defer util.Recover()

		for {
			<-ticker.C
			for _, v := range sm.Keys() {
				li, err := getLogInfoIns(v)
				util.ErrHandler(err)
				if li != nil && li.Sc == sc && time.Since(time.Unix(li.Status[1], 0)) > time.Second*time.Duration(aliveTime) {
					fmt.Println("stop watch: ", v)
					delInfo(v)
				}
			}
		}
	}()
}

func reScanTask(sc *util.SingleConfig) {
	freq := sc.ScanFrequency
	if freq < 1 {
		freq = freqDefault
	}

	go func() {
		defer util.Recover()

		ticker := time.NewTicker(time.Duration(freq) * time.Second)
		for {
			<-ticker.C
			reScan()
		}
	}()
}

func sysInfo() {
	conf := util.GetConfig()
	info := conf.SysInfo

	if info {
		during := conf.SysInfoDuring
		var psInfo gopsinfo.PsInfo
		var d time.Duration
		if during < 1 {
			d = 1
		} else if during == 0 {
			d = 10
		}
		ticker := time.NewTicker(d * time.Second)
		go func() {
			defer util.Recover()
			for {
				<-ticker.C

				psInfo = gopsinfo.GetPsInfo(d)
				sysData, err := json.Marshal(psInfo)
				util.ErrHandler(err)
				doPush(&sysData, systemType, nil)
			}
		}()
	}
}

func takeSnap() {
	conf := util.GetConfig()
	snd := conf.SnapShotDuring
	if snd == 0 {
		snd = snapShotDefault
	}

	ticker := time.NewTicker(time.Duration(snd) * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C

			store := storeState{}
			snap := getSnapPath()
			exist, err := util.PathExist(snap)
			util.ErrHandler(err)
			if !exist {
				err = os.Mkdir(filepath.Dir(snap), os.ModePerm)
			}
			f, err := os.Create(snap)
			for _, v := range sm.Keys() {
				li, err := getLogInfoIns(v)
				util.ErrHandler(err)
				store[v] = li.Status
			}
			d, err := json.Marshal(store)
			util.ErrHandler(err)
			_, err = f.Write(d)
			err = f.Close()
			util.ErrHandler(err)
		}
	}()
}

func debugInfo() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C
			fmt.Println("\n\n\nsm count: ", sm.Count())
			for k, v := range sm.Items() {
				val := v.(logInfo)
				fmt.Println(k, "--->", val)
			}
		}
	}()
}

func getSnapPath() string {
	conf := util.GetConfig()
	snap := conf.SnapshotDir
	if snap == "" {
		snap = filepath.Join(util.GetTempDir(), util.SnapshotDir, util.SnapshotFile)
	}
	return snap
}
