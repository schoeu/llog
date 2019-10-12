package util

import (
	"fmt"
	"reflect"
)

func CombineData(inputVal interface{}) map[string]string {
	fieldVal, ok := inputVal.(map[string]interface{})
	if !ok {
		panic("json unmarshal error.")
	}

	rs := map[string]string{}

	for i, v := range fieldVal {
		t := reflect.TypeOf(v).String()
		if t == "string" {
			rs[i] = v.(string)
		} else if t == "float64" {
			rs[i] = fmt.Sprintf("%.0f", v)
		}
	}

	return rs
}
