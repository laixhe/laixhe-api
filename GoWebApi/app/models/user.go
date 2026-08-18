package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserTable 用户表名
const UserTable = "user"

// UserColumnsNoPassword 查询用户时使用的列 (排除 password)
// 密码哈希只在登录校验时需要, 其它场景读取会把无用的敏感字段拉进内存
const UserColumnsNoPassword = "id,type_id,account,mobile,email,nickname,avatar_url,sex,states,created_at"

// User 用户
type User struct {
	ID        int       `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	TypeId    UserType  `gorm:"column:type_id;type:int;not null;default:0;comment:类型 1普通"`
	Account   string    `gorm:"column:account;type:string;size:120;not null;uniqueIndex;default:'';comment:账号"`
	Mobile    string    `gorm:"column:mobile;type:string;size:120;not null;index;default:'';comment:手机号"`
	Email     string    `gorm:"column:email;type:string;size:120;not null;uniqueIndex;default:'';comment:邮箱"`
	Password  string    `gorm:"column:password;type:string;size:120;not null;default:'';comment:密码"`
	Nickname  string    `gorm:"column:nickname;type:string;size:120;not null;default:'';comment:昵称"`
	AvatarUrl string    `gorm:"column:avatar_url;type:string;size:255;not null;default:'';comment:头像地址"`
	Sex       UserSex   `gorm:"column:sex;type:int;not null;default:0;comment:性别 0未填写 1男 2女"`
	States    UserState `gorm:"column:states;type:int;not null;default:0;comment:状态 0封禁 1正常"`
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:更新时间"`
}

func (m *User) TableName() string {
	return UserTable
}

// CreateUser 在事务中创建用户，同时创建关联的扩展信息和第三方记录
func CreateUser(db *gorm.DB, user *User) error {
	// 事务（返回任何错误都会回滚事务）
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userExtend := &UserExtend{
			Uid: user.ID,
		}
		if err := tx.Create(userExtend).Error; err != nil {
			return err
		}
		userThirdParty := &UserThirdParty{
			Uid: user.ID,
		}
		if err := tx.Create(userThirdParty).Error; err != nil {
			return err
		}
		// 在同一事务中创建用户、扩展信息、第三方关联
		// INSERT INTO `user` (...)
		// INSERT INTO `user_extend` (...)
		// INSERT INTO `user_third_party` (...)
		return nil
	})
}

// UpdateUser 根据非零字段动态更新用户信息，同时更新 updated_at
//
// 注意边界: states 为 0 时无法通过本函数更新为 0 (封禁态);
// 如需显式置 0 需另走 Updates 全量更新。
func UpdateUser(db *gorm.DB, user *User) error {
	if user.ID <= 0 {
		return gorm.ErrPrimaryKeyRequired
	}
	updates := make(map[string]any)
	if user.TypeId > 0 {
		updates["type_id"] = user.TypeId
	}
	if user.Mobile != "" {
		updates["mobile"] = user.Mobile
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Nickname != "" {
		updates["nickname"] = user.Nickname
	}
	if user.AvatarUrl != "" {
		updates["avatar_url"] = user.AvatarUrl
	}
	if user.States > 0 {
		updates["states"] = user.States
	}
	updates["updated_at"] = time.Now()
	// 根据非零字段动态构建 UPDATE 语句
	return db.Model(&User{}).Where("id", user.ID).Updates(updates).Error
}

// ListUser 分页查询用户列表，按 ID 降序
//
// 说明: count(*) 为全表扫描, 数据量达到十万级后每次分页都有一次全表 count;
// 教学规模下可接受, 大表场景建议改用 keyset (游标) 分页或缓存 total。
func ListUser(db *gorm.DB, limit, offset int) ([]User, int, error) {
	var total int64
	var list []User

	if err := db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	// 列表页同样不返回 password
	if err := db.Select(UserColumnsNoPassword).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	// SELECT count(*) FROM `user`
	// SELECT <UserColumnsNoPassword> FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
	return list, int(total), nil
}
