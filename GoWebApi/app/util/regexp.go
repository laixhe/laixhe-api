package util

import (
	"regexp"
)

// matchingPassword 密码正则：字母、数字、_、@、$，长度 >= 6
var matchingPassword = regexp.MustCompile(`^[a-zA-Z0-9_@$]{6,}$`)

// IsPassword 密码是否合法
func IsPassword(password string) bool {
	return matchingPassword.MatchString(password)
}
