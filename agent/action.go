package agent

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/nma/util"
	"github.com/urfave/cli"
)

func StartAction(c *cli.Context) error {
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	conf, err := util.GetConfig(configFile)
	logFiles := conf.LogDir

	if len(logFiles) == 0 {
		logFileDir := util.GetTempDir()
		logFiles = append(logFiles, path.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	logChan := make(chan int)

	// 监控日志收集
	err = fileGlob(logFiles, conf, true)

	// 错误日志收集
	err = fileGlob(conf.ErrLogs, conf, false)

	// 阻塞主goroutines
	<-logChan
	return err
}

func fileGlob(logs []string, conf util.Config, isNormal bool) error {
	for _, v := range logs {
		exist, err := util.PathExist(v)
		util.ErrHandler(err)
		if !exist {
			err = os.Mkdir(v, os.ModePerm)
		}
		if !path.IsAbs(v) {
			v = util.GetAbsPath("", v)
		}
		paths, err := filepath.Glob(v)
		if err != nil {
			return err
		}
		for _, v := range paths {
			go pushLog(v, conf, isNormal)
		}
	}
	return nil
}

func pushLog(logFile string, conf util.Config, isNormal bool) {
	t, err := tail.TailFile(logFile, tail.Config{
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
		Follow: true,
	})

	logServer := util.LogServer
	if conf.LogServer != "" {
		logServer = conf.LogServer
	}

	if isNormal {
		st := time.Now()
		for line := range t.Lines {
			var psInfo gopsinfo.PsInfo
			if !conf.NoSysInfo {
				et := time.Now()
				during := et.Sub(st)
				timeSub := int(during)
				if timeSub < 1 {
					during = time.Microsecond * 1000
				}
				psInfo = gopsinfo.GetPsInfo(during)
				st = et
			}
			var nodeInfo interface{}
			err = json.Unmarshal([]byte(line.Text), &nodeInfo)
			combineRs := util.CombineData(nodeInfo, psInfo, conf.NoSysInfo)
			rs := combineTags(combineRs)
			if logServer != "" {
				PushData(rs, logServer)
			}
		}
		util.ErrHandler(err)
	} else {
		for line := range t.Lines {
			errStruct := map[string]interface{}{
				"errorMsg": line.Text,
			}
			rs := combineTags(errStruct)
			if logServer != "" {
				PushData(rs, logServer)
			}
		}
	}
}

func combineTags(rs map[string]interface{}) map[string]interface{} {
	// 日志签名
	rs["version"] = util.Version
	rs["logId"] = util.UUID()
	rs["type"] = util.AppName
	rs["currentTime"] = time.Now().UnixNano() / 1e6
	return rs
}
