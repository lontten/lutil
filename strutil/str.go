// Package strutil 提供字符串判断与截取工具。
package strutil

import (
	"strings"
	"unicode"
)

// HasText 判断字符串去除空白后是否非空。
func HasText(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0
}

// HasTextP 判断字符串指针去除空白后是否非空；nil 返回 false。
func HasTextP(s *string) bool {
	if s == nil {
		return false
	}
	var s2 = strings.TrimSpace(*s)
	return len(s2) > 0
}

// StrContainsAll 判断 str 是否包含 list 中的全部子串。
func StrContainsAll(str string, list ...string) bool {
	for _, v := range list {
		if !strings.Contains(str, v) {
			return false
		}
	}
	return true
}

// StrContainsAny 判断 str 是否包含 list 中的任一子串。
func StrContainsAny(str string, list ...string) bool {
	for _, v := range list {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

// FirstStrRight 返回 key 在 str 中首次出现位置右侧的子串。
func FirstStrRight(str string, key string) string {
	index := strings.Index(str, key)
	if index == -1 {
		return ""
	}
	start := index + len(key)
	return str[start:]
}

// LastStrRight 返回 key 在 str 中末次出现位置右侧的子串。
func LastStrRight(str string, key string) string {
	index := strings.LastIndex(str, key)
	if index == -1 {
		return ""
	}
	start := index + len(key)
	return str[start:]
}

// LowerFirst 将 str 首字母转为小写。
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
