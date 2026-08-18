package util

import (
	"regexp"
)

// matchingPassword 密码正则：字母、数字、_、@、$，长度 6~64 位
// 长度上限 64: bcrypt 只取前 72 字节, 上限保证任何合法密码都不会被静默截断产生碰撞
var matchingPassword = regexp.MustCompile(`^[a-zA-Z0-9_@$]{6,64}$`)

// IsPassword 密码是否合法
func IsPassword(password string) bool {
	return matchingPassword.MatchString(password)
}
