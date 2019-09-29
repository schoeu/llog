package util

import (
	"fmt"
	"reflect"
)

func CombineData(inputVal interface{}) map[string]interface{} {
	fieldVal, ok := inputVal.(map[string]interface{})
	if !ok {
		panic("json unmarshal error.")
	}

	rs := map[string]interface{}{}

	for i, v := range fieldVal {
		t := reflect.TypeOf(v).String()
		if t == "string" {
			rs[i] = v
		} else if t == "float64" {
			rs[i] = fmt.Sprintf("%.0f", v)
		}
	}

	return rs
}
