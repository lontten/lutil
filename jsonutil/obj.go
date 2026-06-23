// Package jsonutil 提供 JSON 序列化与反序列化便捷函数。
package jsonutil

import "encoding/json"

// ToJsonStr 将 v 序列化为 JSON 字符串，失败时返回 error。
func ToJsonStr(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToJsonStrPanic 将 v 序列化为 JSON 字符串，失败时 panic。
func ToJsonStrPanic(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// ToJsonStrDefault 将 v 序列化为 JSON 字符串，失败时返回空字符串。
func ToJsonStrDefault(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// ToJsonStrP 将 v 序列化为 JSON 字符串指针，失败时返回 error。
func ToJsonStrP(v any) (*string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	s := string(bytes)
	return &s, nil
}

// ToJsonStrPPanic 将 v 序列化为 JSON 字符串指针，失败时 panic。
func ToJsonStrPPanic(v any) *string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	s := string(bytes)
	return &s
}

// ToJsonStrPDefault 将 v 序列化为 JSON 字符串指针，失败时返回指向空字符串的指针。
func ToJsonStrPDefault(v any) *string {
	bytes, err := json.Marshal(v)
	if err != nil {
		s := ""
		return &s
	}
	s := string(bytes)
	return &s
}

// ToObj 将 JSON 字符串反序列化为 T，失败时返回 error。
func ToObj[T any](str string) (T, error) {
	var obj T
	err := json.Unmarshal([]byte(str), &obj)
	return obj, err
}

// ToObjPanic 将 JSON 字符串反序列化为 T，失败时 panic。
func ToObjPanic[T any](str string) T {
	var obj T
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

// ToObjDefault 将 JSON 字符串反序列化为 T，失败时返回零值。
func ToObjDefault[T any](str string) T {
	var obj T
	json.Unmarshal([]byte(str), &obj)
	return obj
}
