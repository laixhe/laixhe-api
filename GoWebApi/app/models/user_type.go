package models

// UserType 用户类型
// 使用定义类型而非别名, 避免与普通 int 混用导致赋值越界
type UserType int

const (
	UserTypeOrdinary UserType = 1 // 普通用户
)

// IsUserTypeValid 判断用户类型是否有效
func IsUserTypeValid(t UserType) bool {
	switch t {
	case UserTypeOrdinary:
		return true
	}
	return false
}

// GetUserTypeText 获取用户类型中文描述
func GetUserTypeText(t UserType) string {
	switch t {
	case UserTypeOrdinary:
		return "普通用户"
	default:
		return ""
	}
}
