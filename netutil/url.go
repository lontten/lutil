package netutil

import (
	"path/filepath"
	"strings"
	"unicode"
)

// unsafeFileNameChars 为 OSS object key / 下载 URL 中不安全、以及本地非法的字符。
// 前半与实战验证过的 reservedChars 一致；后半叠加 Windows 本地非法字符 " < > |。
var unsafeFileNameChars = map[rune]struct{}{
	':': {}, '/': {}, '?': {}, '#': {}, '[': {}, ']': {},
	'@': {}, '!': {}, '$': {}, '&': {}, '\'': {}, '\\': {}, '(': {},
	')': {}, '*': {}, '+': {}, ',': {}, ';': {}, '=': {}, '%': {},
	'"': {}, '<': {}, '>': {}, '|': {},
}

// SafeFileName 将上传原始文件名消毒为可作 OSS object key 的名称，尽量保留可读原文。
// 会替换 URL/OSS 保留字符（含 %）与 Windows 非法字符，避免对象无法下载。
// maxLen == 0 表示不截断；截断时优先保留扩展名。
func SafeFileName(name string, maxLen int) string {
	name = strings.TrimSpace(name)
	name = filepath.Base(name)

	var sb strings.Builder
	sb.Grow(len(name))
	for _, r := range name {
		if unicode.IsSpace(r) || isInvisibleControlCharacter(r) {
			sb.WriteByte('_')
			continue
		}
		if _, ok := unsafeFileNameChars[r]; ok {
			sb.WriteByte('_')
			continue
		}
		sb.WriteRune(r)
	}
	s := sb.String()

	if s == "" || s == "." || s == ".." {
		s = "file"
	}

	if maxLen <= 0 {
		return s
	}
	return truncateKeepingExt(s, maxLen)
}

// SafeURL 将字符串消毒为可作 OSS object key 的文件名。
//
// Deprecated: 请使用 SafeFileName。
func SafeURL(url string, size ...int) string {
	maxLen := 0
	if len(size) > 0 {
		maxLen = size[0]
	}
	return SafeFileName(url, maxLen)
}

// truncateKeepingExt 按 rune 截断，尽量保留扩展名，使总长度 ≤ maxLen。
func truncateKeepingExt(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	ext := filepath.Ext(s)
	extRunes := []rune(ext)
	if len(extRunes) > 0 && len(extRunes) < maxLen {
		keep := maxLen - len(extRunes)
		base := string(runes[:len(runes)-len(extRunes)])
		baseRunes := []rune(base)
		if keep > len(baseRunes) {
			keep = len(baseRunes)
		}
		return string(baseRunes[:keep]) + ext
	}
	return string(runes[:maxLen])
}

// isInvisibleControlCharacter 检查字符是否为不可见的控制字符。
func isInvisibleControlCharacter(r rune) bool {
	return (r >= 0x0000 && r <= 0x001F) || (r >= 0x007F && r <= 0x009F)
}
