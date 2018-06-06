package transfer

import (
	"fmt"
	"reflect"
	"strings"
)

func ToString(stype string, value interface{}) string {
	if nil == value {
		return ""
	}

	var tmp string = ""
	switch stype {
	case "text":
		tmp = fmt.Sprintf(" %s ", reflect.ValueOf(value).Elem())
	case "interger", "numeric":
		tmp = fmt.Sprintf(" %d ", reflect.ValueOf(value).Elem())
	case "real":
		tmp = fmt.Sprintf(" %f ", reflect.ValueOf(value).Elem())
	default:
		tmp = fmt.Sprintf("unknown type:%s.", stype)
	}

	tmp = strings.Trim(tmp, " ")
	return strings.Replace(tmp, "\n", "", -1)
}

func ToInt64(value interface{}) string {

	return ""
}

func toSlice(value interface{}) string {
	return ""
}

func ToMap(value interface{}) string {
	return ""
}
