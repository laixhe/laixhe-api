package entity

import (
	"time"

	"webapi/app/models"
)

// User 用户信息
type User struct {
	Uid int `json:"uid" validate:"required"` // 用户id
	// UserType:
	// * 1 - 普通用户
	TypeId    models.UserType `json:"type_id" validate:"required"`    // 类型
	Account   string          `json:"account" validate:"required"`    // 账号
	Mobile    string          `json:"mobile" validate:"required"`     // 手机号
	Email     string          `json:"email" validate:"required"`      // 邮箱
	Nickname  string          `json:"nickname" validate:"required"`   // 昵称
	AvatarUrl string          `json:"avatar_url" validate:"required"` // 头像地址
	// UserSex:
	// * 0 - 未填写
	// * 1 - 男
	// * 2 - 女
	Sex models.UserSex `json:"sex" validate:"required"` // 性别
	// UserState:
	// * 0 - 禁用
	// * 1 - 正常
	States    models.UserState `json:"states" validate:"required"`     // 状态
	CreatedAt string           `json:"created_at" validate:"required"` // 创建时间
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
	Page     int `query:"page" json:"page" validate:"required"`           // 分页-当前页(默认 1)
	PageSize int `query:"page_size" json:"page_size" validate:"required"` // 分页-每页数量(默认 12)
}

// UserListResponse 响应-获取用户列表
type UserListResponse struct {
	Total    int    `json:"total" validate:"required"`     // 总数
	Page     int    `json:"page" validate:"required"`      // 分页-当前页
	PageSize int    `json:"page_size" validate:"required"` // 分页-每页数量
	List     []User `json:"list" validate:"required"`      // 列表
}

// NewUserFromModel 从 DB 模型转换为响应实体，overrideNickname/overrideAvatarUrl 不为空时覆盖对应字段
func NewUserFromModel(m *models.User, overrideNickname, overrideAvatarUrl string) *User {
	nick, avatar := m.Nickname, m.AvatarUrl
	if overrideNickname != "" {
		nick = overrideNickname
	}
	if overrideAvatarUrl != "" {
		avatar = overrideAvatarUrl
	}
	return &User{
		Uid:       m.ID,
		TypeId:    m.TypeId,
		Account:   m.Account,
		Mobile:    m.Mobile,
		Email:     m.Email,
		Nickname:  nick,
		AvatarUrl: avatar,
		Sex:       m.Sex,
		States:    m.States,
		CreatedAt: m.CreatedAt.Format(time.DateTime),
	}
}
