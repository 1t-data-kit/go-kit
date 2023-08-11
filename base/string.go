package base

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func AssertString(data interface{}) string {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Struct, reflect.Slice, reflect.Map:
		m, _ := json.Marshal(data)
		return string(m)
	}
	return ""
}

func AssertStrings(sep string, data ...interface{}) string {
	_strings := make([]string, 0)
	for _, d := range data {
		_strings = append(_strings, AssertString(d))
	}
	return strings.Join(_strings, sep)
}
