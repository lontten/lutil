package strutil

import (
	"strings"
	"unicode"
)

func HasText(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0
}

func HasTextP(s *string) bool {
	if s == nil {
		return false
	}
	var s2 = strings.TrimSpace(*s)
	return len(s2) > 0
}

func StrContainsAll(str string, list ...string) bool {
	for _, v := range list {
		if !strings.Contains(str, v) {
			return false
		}
	}
	return true
}

func StrContainsAny(str string, list ...string) bool {
	for _, v := range list {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

// 第一个key右边第一个子串
func FirstStrRight(str string, key string) string {
	index := strings.Index(str, key)
	if index == -1 {
		return "" // 未找到子串
	}
	start := index + len(key)
	return str[start:]
}

// 最后一个key右边第一个子串
func LastStrRight(str string, key string) string {
	index := strings.LastIndex(str, key)
	if index == -1 {
		return "" // 未找到子串
	}
	start := index + len(key)
	return str[start:]
}

// 首字母小写
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
