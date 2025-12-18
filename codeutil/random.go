package codeutil

import (
	"crypto/rand"
	"math/big"
)

// RandomStrFromCharset 支持 Unicode 字符集
// charset 里可以含中文、emoji 等任意 Unicode 字符
func RandomStrFromCharset(charset string, length int) string {
	if length <= 0 {
		panic("length must be positive")
	}
	runes := []rune(charset)
	if len(runes) < 2 {
		panic("charset must contain at least 2 runes")
	}

	m := big.NewInt(int64(len(runes)))
	buf := make([]rune, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, m)
		if err != nil {
			panic("crypto/rand unavailable: " + err.Error())
		}
		buf[i] = runes[n.Int64()]
	}
	return string(buf)
}

const (
	AlphaNumCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DigitCharset    = "0123456789"
	HexLowerCharset = "0123456789abcdef"
	HexUpperCharset = "0123456789ABCDEF"

	// Base62（短链、邀请码）
	Base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Base58（比特币地址风格，去掉 0OIl）
	Base58Charset = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// 易读字母数字（去掉 0O1lI）
	FriendlyCharset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// 纯字母（大小写）
	AlphaCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 仅大写字母 + 数字（短信验证码、激活码）
	UpperNumCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 仅小写字母 + 数字（文件名、容器名）
	LowerNumCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// 生成包含大小写字母+数字的随机字符串
func RandomStr(length int) string {
	return RandomStrFromCharset(AlphaNumCharset, length)
}

// 生成包含数字的随机字符串
func RandomNum(length int) string {
	return RandomStrFromCharset(DigitCharset, length)
}

// 生成纯数字验证码
func GenCaptcha(length int) string {
	return RandomNum(length)
}

// RandomHexStr 生成任意长度的小写十六进制字符串
func RandomHexStr(length int) string {
	return RandomStrFromCharset(HexLowerCharset, length)
}

// RandomHexStrUpper 生成任意长度的大写十六进制字符串
func RandomHexStrUpper(length int) string {
	return RandomStrFromCharset(HexUpperCharset, length)
}
