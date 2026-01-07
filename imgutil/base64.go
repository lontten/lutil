package imgutil

import (
	"encoding/base64"
	"errors"
)

// BytesToBase64 将[]byte类型的二进制数据转换为标准Base64字符串
// 参数:
//
//	data - 待转换的二进制数据（如二维码的[]byte）
//
// 返回:
//
//	string - 转换后的Base64字符串
//	error - 转换失败时返回的错误（如数据为空）
func BytesToBase64(data []byte) (string, error) {
	// 空值校验，避免转换空数据
	if len(data) == 0 {
		return "", errors.New("待转换的二进制数据为空")
	}

	// 标准Base64编码（URL安全版用URLEncoding，见下方说明）
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

// 可选：URL安全的Base64转换（若需在URL/JSON中传输，避免+、/等特殊字符）
// BytesToURLSafeBase64 将[]byte转换为URL安全的Base64字符串（替换+为-，/为_）
func BytesToURLSafeBase64(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("待转换的二进制数据为空")
	}

	base64Str := base64.URLEncoding.EncodeToString(data)
	return base64Str, nil
}
