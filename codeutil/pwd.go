package codeutil

// EnPwd 对密码加盐哈希，返回 salt+hash 密文。
func EnPwd(str string) string {
	salt := RandomStr(32)
	return salt + GetSHA256HashCode([]byte(str), salt)
}

// CheckPassword 校验明文密码与 EnPwd 生成的密文是否匹配。
func CheckPassword(pwd string, ciphertext string) bool {
	if len(ciphertext) < 33 {
		return false
	}
	return GetSHA256HashCode([]byte(pwd), ciphertext[:32]) == ciphertext[32:]
}
