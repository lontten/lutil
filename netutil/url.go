package netutil

import (
	"strings"
	"unicode"
)

// 定义RFC 3986规定的保留字符集合
var reservedChars = map[rune]struct{}{
	':': {}, '/': {}, '?': {}, '#': {}, '[': {}, ']': {},
	'@': {}, '!': {}, '$': {}, '&': {}, '\'': {}, '\\': {}, '(': {},
	')': {}, '*': {}, '+': {}, ',': {}, ';': {}, '=': {},
}

// SafeURL 转换为URL安全的格式
func SafeURL(url string, size ...int) string {
	url = strings.TrimSpace(url)
	var sb strings.Builder
	sb.Grow(len(url)) // 预分配内存提高效率

	for _, c := range url {
		if _, ok := reservedChars[c]; ok {
			sb.WriteByte('_')
		} else {
			sb.WriteRune(c)
		}
	}

	s := sb.String()
	s = CleanString(s)
	runes := []rune(s)
	l := 0
	if len(size) > 0 {
		l = size[0]
	}

	// 如果字符总数小于等于15，直接返回原字符串
	if l == 0 || len(runes) <= l {
		return s
	}
	// 否则返回前15个字符
	return string(runes[:l])
}

// CleanString 将字符串中的空格、空字符和不可见字符替换为下划线
func CleanString(s string) string {
	return strings.Map(func(r rune) rune {
		// 检查字符是否为空格、空字符或不可见字符
		if unicode.IsSpace(r) || isInvisibleControlCharacter(r) {
			return '_' // 替换为下划线
		}
		return r // 保留原字符
	}, s)
}

// isInvisibleControlCharacter 检查字符是否为不可见的控制字符
func isInvisibleControlCharacter(r rune) bool {
	// 控制字符的 Unicode 范围是 0x0000 到 0x001F 和 0x007F 到 0x009F
	return (r >= 0x0000 && r <= 0x001F) || (r >= 0x007F && r <= 0x009F)
}
