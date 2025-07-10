package utils

import (
	"regexp"
)

// 密码正则表达式
var passwordMatcher = regexp.MustCompile(`^[a-zA-Z0-9_@$]{6,}$`)

// IsPassword 密码是否合法
func IsPassword(password string) bool {
	return passwordMatcher.MatchString(password)
}
