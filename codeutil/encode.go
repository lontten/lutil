package codeutil

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetSHA256HashCode(message []byte, salt string) string {
	//计算哈希值,返回一个长度为32的数组
	bytes2 := sha256.Sum256([]byte(string(message) + salt))
	//将数组转换成切片,转换成16进制,返回字符串
	hashcode := hex.EncodeToString(bytes2[:])
	return hashcode
}

// base64 加密
func Base64Encode(str string) string {
	input := []byte(str)
	return base64.StdEncoding.EncodeToString(input)
}

// base64 解密
func Base64Decode(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes), err
}
func MD5(str string) string {
	data := []byte(str)

	// 2. 计算MD5哈希（返回16字节的哈希值）
	hashBytes := md5.Sum(data)

	// 3. 转换为32位十六进制字符串（小写）
	return hex.EncodeToString(hashBytes[:])
}
