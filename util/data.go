package util

import (
	"fmt"
	"github.com/schoeu/gopsinfo"
	"reflect"
	"strings"
)

func CombineData(inputVal interface{}, info gopsinfo.PsInfo, noSysInfo bool) map[string]interface{} {
	fieldVal, ok := inputVal.(map[string]interface{})
	if !ok {
		panic("json unmarshal error.")
	}

	rs := map[string]interface{}{}
	if !noSysInfo {
		getType := reflect.TypeOf(info)
		getValue := reflect.ValueOf(info)
		for i := 0; i < getType.NumField(); i++ {
			field := getType.Field(i)
			value := getValue.Field(i).Interface()
			//fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
			if field.Name != "" {
				rs[strings.ToLower(field.Name[:1])+field.Name[1:]] = value
			}
		}
	}

	for i, v := range fieldVal {
		t := reflect.TypeOf(v).String()
		if t == "string" {
			rs[i] = v
		} else if t == "float64" {
			rs[i] = fmt.Sprintf("%.0f", v)
		}
	}

	// 日志签名
	rs["version"] = Version
	rs["logId"] = UUID()
	rs["type"] = "nmAgent"

	return rs
}
