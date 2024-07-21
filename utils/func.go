package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(function any) string {
	funcPath := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	funcPathArray := strings.Split(funcPath, "/")
	funcName := funcPathArray[len(funcPathArray)-1]

	return funcName
}
