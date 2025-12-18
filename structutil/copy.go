package structutil

import "reflect"

// source copy 到 binding  并覆盖
func CopyStruct(source any, binding any, exclude []string) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(source).Elem()  //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	var size = vVal.NumField()
	var i = 0
LOOP:
	for i < size {
		name := vTypeOfT.Field(i).Name
		for j := range exclude {
			if exclude[j] == name {
				i++
				goto LOOP
			}
		}
		if ok := bVal.FieldByName(name).IsValid(); ok {

			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))

		}
		i++
	}
}
