package codeutil

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt（DefaultCost）对密码哈希。
// 密码存储请优先使用此函数。
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword 判断 password 是否与 HashPassword 生成的 bcrypt 哈希匹配。
func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// EnPwd 对密码加盐并用 SHA-256 哈希，返回 salt+hash 密文。
//
// Deprecated: EnPwd 使用快速 SHA-256，不适合密码存储。
// 请改用 HashPassword。
func EnPwd(str string) string {
	salt := RandomStr(32)
	return salt + GetSHA256HashCode([]byte(str), salt)
}

// CheckPassword 判断 pwd 是否与 EnPwd 密文匹配。
//
// Deprecated: CheckPassword 仅校验 EnPwd（SHA-256）哈希。
// 请改用 VerifyPassword 与 HashPassword。
func CheckPassword(pwd string, ciphertext string) bool {
	if len(ciphertext) < 33 {
		return false
	}
	return GetSHA256HashCode([]byte(pwd), ciphertext[:32]) == ciphertext[32:]
}
