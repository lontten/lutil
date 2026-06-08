package jsonutil

import "encoding/json"

func ToJsonStr(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ToJsonStrPanic(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ToJsonStrDefault(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ToJsonStrP(v any) (*string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	s := string(bytes)
	return &s, nil
}

func ToJsonStrPPanic(v any) *string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s := string(bytes)
	return &s
}

func ToObj[T any](str string) (T, error) {
	var obj T
	err := json.Unmarshal([]byte(str), &obj)
	return obj, err
}

func ToObjPanic[T any](str string) T {
	var obj T
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func ToObjDefault[T any](str string) T {
	var obj T
	json.Unmarshal([]byte(str), &obj)
	return obj
}
