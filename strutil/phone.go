package strutil

import (
	"regexp"
)

func CheckPhoneAll(phoneNumber string) bool {
	// 手机号
	phone := CheckPhone(phoneNumber)
	if phone {
		return true
	}

	//广义范围的电话号码
	hyphen := isValidNumberWithHyphen(phoneNumber)
	if hyphen {
		return true
	}
	return false
}

func CheckPhone(phoneNumber string) bool {
	// 手机号
	pattern := "^1[3-9][0-9]{9}$"
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}

func CheckLandline(phoneNumber string) bool {
	// 固话 区号（2-3位） + 连字符 + 本地号码（7-8位） + 连字符 + 分机号（2-4位）
	pattern := `^0\d{2,3}-\d{7,8}(-\d{2,4})?$`
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}

// 判断是否符合数字加横杠格式
func isValidNumberWithHyphen(s string) bool {
	// 长度在4到30位之间
	if len(s) < 4 || len(s) > 22 {
		return false
	}
	// 只包含数字和横杠，且横杠不能连续出现，也不能出现在开头或结尾
	matched, _ := regexp.MatchString(`^\d+(?:-\d+)*$`, s)
	return matched
}
