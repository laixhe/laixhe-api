package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserTable 用户表名
const UserTable = "user"

// User 用户
type User struct {
	ID        int       `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	TypeId    UserType  `gorm:"column:type_id;type:int;not null;default:0;comment:类型 1普通"`
	Account   string    `gorm:"column:account;type:string;size:120;not null;index;default:'';comment:账号"`
	Mobile    string    `gorm:"column:mobile;type:string;size:120;not null;index;default:'';comment:手机号"`
	Email     string    `gorm:"column:email;type:string;size:120;not null;index;default:'';comment:邮箱"`
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

func (m *User) Get(db *gorm.DB, user User) error {
	wheres := make(map[string]any)
	if user.ID > 0 {
		wheres["id"] = user.ID
	}
	if user.Email != "" {
		wheres["email"] = user.Email
	}
	if user.Mobile != "" {
		wheres["mobile"] = user.Mobile
	}
	if user.Nickname != "" {
		wheres["nickname"] = user.Nickname
	}
	if len(wheres) == 0 {
		return gorm.ErrRecordNotFound
	}
	// SELECT * FROM `user` WHERE `id` = ? ORDER BY id LIMIT 1
	// SELECT * FROM `user` WHERE `email` = ? ORDER BY id LIMIT 1
	// SELECT * FROM `user` WHERE `mobile` = ? ORDER BY id LIMIT 1
	// SELECT * FROM `user` WHERE `nickname` = ? ORDER BY id LIMIT 1
	return db.Where(wheres).First(m).Error
}

func (m *User) Create(db *gorm.DB) error {
	// 事务（返回任何错误都会回滚事务）
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		userExtend := &UserExtend{
			Uid: m.ID,
		}
		if err := tx.Create(userExtend).Error; err != nil {
			return err
		}
		userThirdParty := &UserThirdParty{
			Uid: m.ID,
		}
		if err := tx.Create(userThirdParty).Error; err != nil {
			return err
		}
		// INSERT INTO `user` (`type_id`,`mobile`,`email`,`password`,`nickname`,`avatar_url`,`states`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)
		// INSERT INTO `user_extend` (`uid`,`birthday`,`height`,`weight`) VALUES (?,?,?)
		// INSERT INTO `user_third_party` (`uid`,`wechat_unionid`,`wechat_openid`) VALUES (?,?,?)
		return nil
	})
}

func (m *User) Update(db *gorm.DB) error {
	if m.ID <= 0 {
		return gorm.ErrPrimaryKeyRequired
	}
	updates := make(map[string]any)
	if m.TypeId > 0 {
		updates["type_id"] = m.TypeId
	}
	if m.Mobile != "" {
		updates["mobile"] = m.Mobile
	}
	if m.Email != "" {
		updates["email"] = m.Email
	}
	if m.Password != "" {
		updates["password"] = m.Password
	}
	if m.Nickname != "" {
		updates["nickname"] = m.Nickname
	}
	if m.AvatarUrl != "" {
		updates["avatar_url"] = m.AvatarUrl
	}
	if m.States > 0 {
		updates["states"] = m.States
	}
	// UPDATE `user` SET `nickname`=?,`updated_at`=? WHERE `id` = ?
	return db.Model(m).Updates(updates).Error
}

func (m *User) List(db *gorm.DB, limit, offset int) ([]User, int, error) {
	var total int64
	var list []User

	if err := db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	if err := db.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(limit).
		Offset(offset).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	// SELECT count(*) FROM `user`
	// SELECT * FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
	return list, int(total), nil
}
