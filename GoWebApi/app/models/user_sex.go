package models

// UserSex 用户性别
// 使用定义类型而非别名, 避免与普通 int 混用导致赋值越界
type UserSex int

const (
	UserSexUnknown UserSex = 0 // 未填写
	UserSexMale    UserSex = 1 // 男
	UserSexFemale  UserSex = 2 // 女
)

// IsUserSexValid 判断性别值是否有效（仅男/女为有效值）
func IsUserSexValid(s UserSex) bool {
	switch s {
	case UserSexMale, UserSexFemale:
		return true
	}
	return false
}

// GetUserSexText 获取性别中文描述
func GetUserSexText(s UserSex) string {
	switch s {
	case UserSexMale:
		return "男"
	case UserSexFemale:
		return "女"
	default:
		return ""
	}
}
