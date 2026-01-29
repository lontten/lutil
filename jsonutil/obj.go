package jsonutil

import "encoding/json"

func ToJsonStr(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ToJsonStrP(v any) *string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s := string(bytes)
	return &s
}

func ToObj[T any](str string) T {
	var obj T
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func ToObjList[T any](str string) []T {
	var obj []T
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		panic(err)
	}
	return obj
}
