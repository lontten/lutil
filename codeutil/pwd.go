package codeutil

// 密码加密
func EnPwd(str string) string {
	salt := RandomStr(32)
	return salt + GetSHA256HashCode([]byte(str), salt)
}

// 密码校验
func CheckPassword(pwd string, ciphertext string) bool {
	if len(ciphertext) < 33 {
		return false
	}
	return GetSHA256HashCode([]byte(pwd), ciphertext[:32]) == ciphertext
}
