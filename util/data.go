package util

import (
	"fmt"
	"github.com/schoeu/gopsinfo"
	"reflect"
)

func CombineData(inputVal interface{}, info gopsinfo.PsInfo) map[string]interface{} {
	fieldVal, ok := inputVal.(map[string]interface{})
	if !ok {
		panic("json unmarshal error.")
	}

	getType := reflect.TypeOf(info)
	getValue := reflect.ValueOf(info)
	rs := map[string]interface{}{}
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		//fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		if field.Name != "" {
			rs[field.Name] = value
		}
	}

	for i, v := range fieldVal {
		rs[i] = fmt.Sprintf("%.0f", v)
	}

	rs["version"] = Version
	rs["logId"] = UUID()
	rs["type"] = "nmAgent"

	return rs
}
