package codeutil

import (
	"time"
)

// RandomTimeID32
// 格式 yyyymmddHHMMSS + 18位随机数,有字母
func RandomTimeID32() string {
	return RandomTimeID(32)
}

// RandomTimeID
// 格式 yyyymmddHHMMSS + n位随机数,有字母
func RandomTimeID(length int) string {
	return RandomTimeBaseID(UpperNumCharset, length)
}

// RandomTimeNumberID32
// 格式 yyyymmddHHMMSS + 18位随机数字
func RandomTimeNumberID32() string {
	return RandomTimeNumberID(32)
}

// RandomTimeNumberID
// 格式 yyyymmddHHMMSS + n位随机数字
func RandomTimeNumberID(length int) string {
	return RandomTimeBaseID(DigitCharset, length)
}

func RandomTimeBaseID(charset string, length int) string {
	timestamp := time.Now().Format("20060102150405") // 14字符
	randomLength := length - len(timestamp)
	if randomLength <= 0 {
		return timestamp[:length]
	}
	return timestamp + RandomStrFromCharset(charset, randomLength)
}
