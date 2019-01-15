package ds

import (
	"reflect"
	"strings"
)

type ListElement struct {
	Next *ListElement
	Prev *ListElement
	Value interface{}
}

func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("ToInterfaceSlice given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func RemoveFromSlice(slice []interface{}, v interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}

	idx := -1
	for i, e := range slice {
		if e == v {
			idx = i
		}
	}
	if idx < 0 {
		return slice
	}

	last := slice[idx+1:]
	slice = append(slice[:idx], last...)
	return slice
}

type ErrorList []error

func (e ErrorList) Error() string {
	builder := strings.Builder{}
	for i, err := range e {
		builder.WriteString(err.Error())
		if i != len(e) - 1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
