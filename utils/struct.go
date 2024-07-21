package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func GetKeysFromStruct(v any) (keys []string, err error) {
	val := reflect.ValueOf(v)
	vType := reflect.TypeOf(v)

	if vType.Kind() != reflect.Struct {
		return keys, errors.New("value is not a struct")
	}

	for i := range val.NumField() {
		keys = append(keys, vType.Field(i).Name)
	}

	return keys, nil
}

func StructToMap(v any) (mappedStruct map[string]any, err error) {
	// Check if v is already a map and return it if true
	if mappedStruct, ok := v.(map[string]any); ok {
		return mappedStruct, nil
	}

	mappedStruct = make(map[string]any)

	val := reflect.ValueOf(v)
	vType := reflect.TypeOf(v)

	// Check if v is a pointer and dereference it
	if vType.Kind() == reflect.Ptr {
		val = val.Elem()
		vType = vType.Elem()
	}

	if vType.Kind() != reflect.Struct {
		return mappedStruct, errors.New("value is not a struct, " + vType.Elem().String())
	}

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i).Elem()

		if value != reflect.ValueOf(nil) {
			key := vType.Field(i).Name

			key = strings.ToLower(string(key[0])) + key[1:]

			mappedStruct[key] = value
		}
	}

	return mappedStruct, nil
}

func PrintStructType[T any](val T) string {
	t := reflect.TypeOf(val)
	str := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		str += fmt.Sprintf("  %s (%s)\n", field.Name, field.Type)
	}

	return str
}
