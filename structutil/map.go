package structutil

import (
	"fmt"
	"reflect"
)

func Struct2StringMap(s any) map[string]string {
	val := reflect.ValueOf(s)

	// 解引用指针（支持结构体指针输入）
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	relType := val.Type()
	numField := relType.NumField()
	m := make(map[string]string, numField)

	for i := 0; i < numField; i++ {
		field := relType.Field(i)
		fieldVal := val.Field(i)

		// 转换字段值为字符串（支持多种类型）
		var strVal string
		switch fieldVal.Kind() {
		case reflect.String:
			strVal = fieldVal.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strVal = fmt.Sprintf("%d", fieldVal.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			strVal = fmt.Sprintf("%d", fieldVal.Uint())
		case reflect.Float32, reflect.Float64:
			strVal = fmt.Sprintf("%v", fieldVal.Float())
		case reflect.Bool:
			strVal = fmt.Sprintf("%t", fieldVal.Bool())
		default:
			// 处理其他未覆盖的类型（如切片、嵌套结构体等，可根据需求扩展）
			strVal = fmt.Sprintf("%v", fieldVal.Interface())
		}

		m[field.Name] = strVal
	}

	return m
}
