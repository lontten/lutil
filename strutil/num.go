package strutil

import (
	"strings"
	"unicode"
)

// IsLooseNumber 判断字符串是否为宽松格式的数字
func IsLooseNumber(str string) bool {
	// 移除允许的符号
	str = strings.Map(func(r rune) rune {
		switch r {
		case ' ', '+', '-', '*', '/', '.', ',', '_', '%':
			return -1 // 删除这些字符
		}
		return r
	}, str)

	// 空字符串不是数字
	if len(str) == 0 {
		return false
	}

	// 检查是否全是数字
	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// IsNotLooseNumber 判断字符串是否不是宽松格式的数字
// 与 IsLooseNumber 逻辑相反，用于保持向后兼容
func IsNotLooseNumber(str string) bool {
	return !IsLooseNumber(str)
}
