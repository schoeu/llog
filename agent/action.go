package agent

import (
	"encoding/json"
	"fmt"
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

func StartAction(c *cli.Context) {
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	conf, err := util.GetConfig(configFile)
	logFiles := conf.LogDir

	var ch = make(chan int)

	if len(logFiles) == 0 {
		logFileDir := util.GetTempDir()
		logFiles = append(logFiles, path.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	// 监控日志收集
	err = fileGlob(logFiles, conf, true)

	// 错误日志收集
	err = fileGlob(conf.ErrLogs, conf, false)

	util.ErrHandler(err)

	<-ch
}

func fileGlob(logs []string, conf util.Config, isNormal bool) error {
	for _, v := range logs {
		exist, err := util.PathExist(v)
		util.ErrHandler(err)
		if !exist {
			err = os.Mkdir(path.Dir(v), os.ModePerm)
		}
		if !path.IsAbs(v) {
			v = util.GetAbsPath("", v)
		}
		paths, err := filepath.Glob(v)
		if err != nil {
			return err
		}
		for _, v := range paths {
			excludeFiles := conf.ExcludeFiles
			if len(excludeFiles) > 0 && !util.IsInclude(v, conf.ExcludeFiles) {
				continue
			}
			go pushLog(v, conf, isNormal)
		}
	}
	return nil
}

func pushLog(logFile string, conf util.Config, isNormal bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

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

	st := time.Now()
	var rs map[string]interface{}
	include, exclude := conf.Include, conf.Exclude
	for line := range t.Lines {
		if len(include) > 0 && !util.IsInclude(line.Text, include) {
			continue
		}
		if len(exclude) > 0 && util.IsInclude(line.Text, exclude) {
			continue
		}

		if isNormal {
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

			rs = util.CombineData(nodeInfo, psInfo, conf.NoSysInfo)
		} else {
			rs = map[string]interface{}{
				"@message": line.Text,
			}
		}
		rs := combineTags(rs)
		if logServer != "" {
			PushData(rs, logServer)
		}
	}
	util.ErrHandler(err)
}

func combineTags(rs map[string]interface{}) map[string]interface{} {
	// 日志签名
	rs["@version"] = util.Version
	rs["@logId"] = util.UUID()
	rs["@type"] = util.AppName
	rs["@timestamps"] = time.Now().UnixNano() / 1e6
	return rs
}

func StopAction(c *cli.Context) {
	fmt.Println("stop")
}
