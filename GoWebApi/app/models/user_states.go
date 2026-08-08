package models

// UserState 用户账号状态（0=封禁 1=正常）
// 使用定义类型而非别名, 避免与普通 int 混用导致赋值越界
type UserState int

const (
	UserStateBanned UserState = 0 // 禁用
	UserStateNormal UserState = 1 // 正常
)

// IsUserStateValid 判断用户状态值是否有效
func IsUserStateValid(s UserState) bool {
	switch s {
	case UserStateBanned, UserStateNormal:
		return true
	}
	return false
}

// GetUserStateText 获取用户状态中文描述
func GetUserStateText(s UserState) string {
	switch s {
	case UserStateBanned:
		return "禁用"
	case UserStateNormal:
		return "正常"
	default:
		return ""
	}
}
