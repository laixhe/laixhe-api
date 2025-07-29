package entity

import "webapi/app/model"

// User 用户信息
type User struct {
	Uid int `json:"uid" validate:"required"` // 用户id
	// UserType:
	// * 1 - 普通用户
	TypeId    model.UserType `json:"type_id" validate:"required"`    // 类型
	Account   string         `json:"account" validate:"required"`    // 账号
	Mobile    string         `json:"mobile" validate:"required"`     // 手机号
	Email     string         `json:"email" validate:"required"`      // 邮箱
	Nickname  string         `json:"nickname" validate:"required"`   // 昵称
	AvatarUrl string         `json:"avatar_url" validate:"required"` // 头像地址
	// UserSex:
	// * 0 - 未填写
	// * 1 - 男
	// * 2 - 女
	Sex model.UserSex `json:"sex" validate:"required"` // 性别
	// UserState:
	// * 0 - 禁用
	// * 1 - 正常
	States model.UserState `json:"states" validate:"required"` // 状态
}

// UserUpdateRequest 请求-更新用户信息
type UserUpdateRequest struct {
	Uid       int    `json:"-"`                              // 用户id
	Nickname  string `json:"nickname" validate:"required"`   // 昵称
	AvatarUrl string `json:"avatar_url" validate:"required"` // 头像地址
}

// UserInfoRequest 请求-获取用户信息
type UserInfoRequest struct {
	Uid int `query:"uid" validate:"required"` // 用户id
}

// UserListRequest 请求-获取用户列表
type UserListRequest struct {
	PageSize int `query:"page_size" validate:"required"` // 每页数量（分页）
	OffsetId int `query:"offset_id" validate:"required"` // 偏移id（分页）
}

// UserListResponse 响应-获取用户列表
type UserListResponse struct {
	Total int    `json:"total" validate:"required"` // 总数
	List  []User `json:"list" validate:"required"`  // 列表
}
