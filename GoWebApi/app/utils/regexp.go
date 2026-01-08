package utils

import (
	"regexp"
)

// MatchingPassword 定义匹配密码正则表达式
var MatchingPassword = regexp.MustCompile(`^[a-zA-Z0-9_@$]{6,}$`)

// IsPassword 密码是否合法
func IsPassword(password string) bool {
	return MatchingPassword.MatchString(password)
}
