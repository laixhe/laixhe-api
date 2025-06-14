package models

import (
	"context"
	"time"
)

const UserTable = "user"

// User 用户
type User struct {
	ID        int       `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	Typeid    int       `gorm:"column:typeid;type:int;not null;default:0;comment:类型"`
	Mobile    string    `gorm:"column:mobile;type:string;size:100;not null;default:'';comment:手机号"`
	Email     string    `gorm:"column:email;type:string;size:100;not null;default:'';comment:邮箱"`
	Password  string    `gorm:"column:password;type:string;size:120;not null;default:'';comment:密码"`
	Nickname  string    `gorm:"column:nickname;type:string;size:100;not null;default:'';comment:昵称"`
	AvatarUrl string    `gorm:"column:avatar_url;type:string;size:255;not null;default:'';comment:头像地址"`
	States    int       `gorm:"column:states;type:int;not null;default:0;comment:状态 1正常 2封禁"`
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:更新时间"`
}

func (*User) TableName() string {
	return UserTable
}

func (m *Model) GetUserByNickname(ctx context.Context, nickname string) (*User, error) {
	var user User
	if err := m.orm().WithContext(ctx).
		Where("nickname", nickname).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
