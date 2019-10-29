package agent

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

var apiServer, name string

const errorType = "error"
const normalType = "normal"
const systemType = "system"

var syncMapError = errors.New("sync map error")

func fileGlob(sc *util.SingleConfig) {
	allLogs := sc.LogDir
	if len(allLogs) == 0 {
		logFileDir := util.GetTempDir()
		allLogs = append(allLogs, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}

	// allLogs: - /var/logs/**/*.log
	for _, v := range allLogs {
		v = pathPreProcess(v)
		// paths: ["/var/logs/1.log","/var/logs/2.log"]
		paths, err := filepath.Glob(v)
		util.ErrHandler(err)
		// update file state.
		initState(paths, sc)
	}
}

func pathPreProcess(p string) string {
	exist, err := util.PathExist(p)
	util.ErrHandler(err)
	if !exist {
		err = os.Mkdir(filepath.Dir(p), os.ModePerm)
	}
	if !filepath.IsAbs(p) {
		p = util.GetAbsPath("", p)
	}
	return p
}

func lineFilter(k string) func(*[]byte) {
	fi := getLogInfoIns(k)
	sc := fi.sc

	include, exclude, multiline := sc.Include, sc.Exclude, sc.Multiline.Pattern
	confMaxByte, maxLines := sc.MaxBytes, sc.Multiline.MaxLines

	if maxLines == 0 {
		maxLines = 10
	}

	return func(l *[]byte) {
		line := *l
		// multiple mode
		if multiline != "" {
			buf := fi.data
			// multiple head line
			if util.IsInclude(line, []string{multiline}) {
				if buf.Len() > 0 {
					fmt.Println(buf.String())
					ok, rs := filter(include, exclude, buf.Bytes(), confMaxByte)
					if ok {
						return
					}
					fmt.Println(k, "-->", buf.String())
					doPush(rs, errorType)

					sm.Set(k, logInfo{
						data:    bytes.Buffer{},
						sc:      fi.sc,
						status:  fi.status,
						fileIns: fi.fileIns,
					})
					fi.lineCount = 0
				}
			}
			fi.lineCount++
			// 匹配多行其他内容
			if fi.lineCount < maxLines {
				//logContent.Write(line)
				buf.Write(line)
				sm.Set(k, logInfo{
					data:    buf,
					sc:      fi.sc,
					status:  fi.status,
					fileIns: fi.fileIns,
					//lineCount: fi.lineCount,
				})
			}
		} else {
			ok, rs := filter(include, exclude, line, confMaxByte)
			if ok {
				return
			}
			doPush(rs, normalType)
		}
	}
}

func filter(include, exclude []string, line []byte, max int) (bool, *[]byte) {
	if len(include) > 0 && !util.IsInclude(line, include) {
		return true, nil
	}
	if len(exclude) > 0 && util.IsInclude(line, exclude) {
		return true, nil
	}
	if max != 0 && len(line) > max {
		line = line[:max]
	}
	return false, &line
}

func doPush(text *[]byte, types string) {
	// 日志签名
	var rs = logStruct{
		"@message":    string(*text),
		"@version":    util.Version,
		"@logId":      util.UUID(),
		"@timestamps": strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		"@types":      types,
		"@name":       name,
	}

	if apiServer != "" {
		go apiPush(&rs, apiServer)
	}

	if indexServer != nil {
		go esPush(&rs)
	}
}

func getLogInfoIns(p string) *logInfo {
	logContent, ok := sm.Get(p)
	if !ok {
		util.ErrHandler(syncMapError)
	}
	li, ok := logContent.(logInfo)
	return &li
}
