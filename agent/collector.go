package agent

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/schoeu/llog/config"
	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

var apiServer, name string
var buf = bytes.Buffer{}
var count = 0
var maxLinesDefault = 10
var snapShotDefault = 5
var json = jsoniter.ConfigCompatibleWithStandardLibrary

const errorType = "error"
const normalType = "normal"
const systemType = "system"

func fileGlob(sc config.SingleConfig) {
	allLogs := sc.LogDir
	if len(allLogs) == 0 {
		logFileDir := util.GetTempDir()
		allLogs = append(allLogs, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}
	// allLogs: - /var/logs/**/*.log
	for _, v := range allLogs {
		v = pathPreProcess(v)
		// paths: ["/var/logs/1.log","/var/logs/2.log"]
		p, err := filepath.Glob(v)
		util.ErrHandler(err)
		if len(p) > 0 {
			initState(p, sc)
		}
	}
}

func recoverState() {
	snap := getSnapPath()
	if snap != "" {
		d, err := ioutil.ReadFile(snap)
		util.ErrHandler(err)
		ss := storeState{}
		err = json.Unmarshal(d, &ss)
		if err != nil {
			fmt.Println("[LLOG] recover file state error")
			return
		}

		for k, v := range ss {
			li, err := getLogInfoIns(k)
			if li != nil && err == nil {
				sm.Set(k, logInfo{
					Sc:      li.Sc,
					Status:  v,
					FileIns: li.FileIns,
				})
			}
		}
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
	fi, err := getLogInfoIns(k)
	util.ErrHandler(err)
	if fi != nil {
		sc := fi.Sc

		include, exclude, multiline := sc.Include, sc.Exclude, sc.Multiline.Pattern
		confMaxByte, maxLines, fields := sc.MaxBytes, sc.Multiline.MaxLines, sc.Fields

		if maxLines == 0 {
			maxLines = maxLinesDefault
		}

		return func(l *[]byte) {
			line := *l
			// fmt.Println(string(line))
			// multiple mode
			if multiline != "" {
				// multiple head line
				if util.IsInclude(line, []string{multiline}) {
					if buf.Len() > 0 {
						ok, rs := filter(include, exclude, buf.Bytes(), confMaxByte)
						if ok {
							return
						}
						doPush(rs, errorType, fields)
						count = 0
						buf = bytes.Buffer{}
					}
				}
				count++
				if count < maxLines {
					buf.Write(line)
				}
			} else {
				ok, rs := filter(include, exclude, line, confMaxByte)
				if ok {
					return
				}
				doPush(rs, normalType, fields)
			}
		}
	}
	return nil
}
