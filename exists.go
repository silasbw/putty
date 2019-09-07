package main

import (
	"reflect"
	"strconv"
)

func Exists(object interface{}, path []string) bool {
	value := reflect.Indirect(reflect.ValueOf(object))
	for _, piece := range path {
		if value.Kind() == reflect.Struct {
			value = value.FieldByName(piece)
			if !value.IsValid() {
				return false
			}
		} else if value.Kind() == reflect.Slice {
			index, err := strconv.ParseInt(piece, 10, 64)
			if err != nil {
				return false
			}
			if int(index) >= value.Len() {
				return false
			}
			value = value.Index(int(index))
		} else {
			return false
		}
	}
	return true
}
