package entity

import "webapi/app/model"

// User 用户信息
type User struct {
	Uid       int            `json:"uid"`        // 用户id
	TypeId    model.UserType `json:"type_id"`    // 类型id
	Account   string         `json:"account"`    // 账号
	Mobile    string         `json:"mobile"`     // 手机号
	Email     string         `json:"email"`      // 邮箱
	Nickname  string         `json:"nickname"`   // 昵称
	AvatarUrl string         `json:"avatar_url"` // 头像地址
	Sex       model.Sex      `json:"sex"`        // 性别
	States    model.State    `json:"states"`     // 状态
}

// UserUpdateRequest 请求-更新用户信息
type UserUpdateRequest struct {
	Uid       int    `json:"-"`          // 用户id
	Nickname  string `json:"nickname"`   // 昵称
	AvatarUrl string `json:"avatar_url"` // 头像地址
}

// UserInfoRequest 请求-获取用户信息
type UserInfoRequest struct {
	Uid int `query:"uid"` // 用户id
}

// UserListRequest 请求-获取用户列表
type UserListRequest struct {
	PageSize int `query:"page_size"` // 每页数量（分页）
	OffsetId int `query:"offset_id"` // 偏移id（分页）
}

// UserListResponse 响应-获取用户列表
type UserListResponse struct {
	Total int    `json:"total"` // 总数
	List  []User `json:"list"`  // 列表
}
