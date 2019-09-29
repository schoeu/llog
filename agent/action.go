package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
	"github.com/urfave/cli"
)

func StartAction(c *cli.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	configFile := util.GetAbsPath(util.GetCwd(), c.Args().First())
	conf, err := util.GetConfig(configFile)
	util.ErrHandler(err)

	logFiles := conf.LogDir

	var ch = make(chan int)

	if len(logFiles) == 0 {
		logFileDir := util.GetTempDir()
		logFiles = append(logFiles, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	// 监控日志收集
	err = fileGlob(logFiles, conf)

	util.ErrHandler(err)

	<-ch
}

func fileGlob(logs []string, conf util.Config) error {
	for _, v := range logs {
		exist, err := util.PathExist(v)
		util.ErrHandler(err)
		if !exist {
			err = os.Mkdir(filepath.Dir(v), os.ModePerm)
		}
		if !filepath.IsAbs(v) {
			v = util.GetAbsPath("", v)
		}
		paths, err := filepath.Glob(v)
		if err != nil {
			return err
		}
		for _, v := range paths {
			excludeFiles := conf.ExcludeFiles
			if len(excludeFiles) > 0 && util.IsInclude(v, excludeFiles) {
				continue
			}
			go pushLog(v, conf)
		}
	}
	return nil
}

func pushLog(logFile string, conf util.Config) {
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
	util.ErrHandler(err)

	st := time.Now()
	var logContent bytes.Buffer

	include, exclude, apiServer, multiline := conf.Include, conf.Exclude, conf.ApiServer, conf.Multiline.Pattern
	SysInfo, confMaxByte, maxLines := conf.SysInfo, conf.MaxBytes, conf.Multiline.MaxLines
	for line := range t.Lines {
		text := line.Text
		if len(include) > 0 && !util.IsInclude(text, include) {
			continue
		}
		if len(exclude) > 0 && util.IsInclude(text, exclude) {
			continue
		}
		// 多行模式
		if multiline != "" {
			// 匹配开始头
			if util.IsInclude(text, []string{multiline}) {
				if logContent.Len() > 0 {
					doPush(SysInfo, st, logContent.Bytes(), apiServer, confMaxByte)
					logContent = bytes.Buffer{}
				}
			}
			// 匹配多行其他内容
			if maxLines != 0 && logContent.Len() < maxLines {
				logContent.WriteString(text)
				continue
			}
		} else {
			doPush(SysInfo, st, []byte(text), apiServer, confMaxByte)
		}
	}
	util.ErrHandler(err)
}

func doPush(SysInfo bool, st time.Time, text []byte, apiServer string, confMaxByte int) {
	var rs map[string]interface{}

	if confMaxByte != 0 && len(text) > confMaxByte {
		text = text[:confMaxByte]
	}

	if SysInfo {
		var psInfo gopsinfo.PsInfo
		et := time.Now()
		during := et.Sub(st)
		timeSub := int(during)
		if timeSub < 1 {
			during = time.Microsecond * 1000
		}
		psInfo = gopsinfo.GetPsInfo(during)
		st = et

		sysInfo, err := json.Marshal(psInfo)
		util.ErrHandler(err)
		rs = util.CombineData(map[string]interface{}{
			"@sysInfo": string(sysInfo),
			"@message": string(text),
		})
	} else {
		rs = map[string]interface{}{
			"@message": string(text),
		}
	}

	if apiServer != "" {
		go PushData(combineTags(rs), apiServer)
	}
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
