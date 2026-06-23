// Package codeutil 提供编码、哈希、随机字符串与密码工具。
package codeutil

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GetSHA256HashCode 计算 message 与 salt 拼接后的 SHA256 十六进制摘要。
func GetSHA256HashCode(message []byte, salt string) string {
	bytes2 := sha256.Sum256([]byte(string(message) + salt))
	return hex.EncodeToString(bytes2[:])
}

// Base64Encode 对字符串进行标准 Base64 编码。
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode 对 Base64 字符串解码。
func Base64Decode(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes), err
}

// MD5 计算字符串的 MD5 十六进制摘要（小写）。
func MD5(str string) string {
	hashBytes := md5.Sum([]byte(str))
	return hex.EncodeToString(hashBytes[:])
}
