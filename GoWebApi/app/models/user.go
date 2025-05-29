package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/laixhe/gonet"
)

// UserTable 表名
const UserTable = "user"

// User 用户表
type User struct {
	ID        int64          `gorm:"column:id;type:int unsigned;not null;AUTO_INCREMENT;primaryKey;comment:用户ID自增"`
	Password  string         `gorm:"column:password;type:string;size:120;not null;default:'';comment:密码"`
	Email     string         `gorm:"column:email;type:string;size:100;not null;default:'';comment:邮箱"`
	Uname     string         `gorm:"column:uname;type:string;size:100;not null;default:'';comment:用户名"`
	Age       int32          `gorm:"column:age;type:tinyint unsigned;not null;default:0;comment:年龄"`
	Score     float64        `gorm:"column:score;type:decimal(10,2) unsigned;not null;default:0.00;comment:分数"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间"`
}

func (*User) TableName() string {
	return UserTable
}

func (u *User) Create(user *User) error {
	// INSERT INTO `user` (`password`,`email`,`uname`,`age`,`score`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?)
	return gonet.Gorm().Client().Create(user).Error
}

func (u *User) FirstEmail(email string) (User, error) {
	var user User
	// SELECT * FROM `user` WHERE email = ? AND `deleted_at` IS NULL ORDER BY `id` LIMIT 1
	err := gonet.Gorm().Client().Where("email", email).First(&user).Error
	return user, err
}

func (u *User) FirstUname(uname string) (User, error) {
	var user User
	// SELECT * FROM `user` WHERE `uname` = ? AND `deleted_at` IS NULL ORDER BY `id` LIMIT 1
	err := gonet.Gorm().Client().Where("uname", uname).First(&user).Error
	return user, err
}

func (u *User) FirstID(uid int64) (User, error) {
	var user User
	// SELECT * FROM `user` WHERE id = ? AND `deleted_at` IS NULL ORDER BY `id` LIMIT 1
	// err := gonet.Gorm().Where("id = ?", uid).First(&user).Error

	// SELECT * FROM `user` WHERE id = ? AND `deleted_at` IS NULL LIMIT 1
	err := gonet.Gorm().Client().Where("id", uid).Take(&user).Error
	return user, err
}

func (u *User) List(size, page int) ([]User, int64, error) {
	var users []User
	// SELECT * FROM `user` WHERE `deleted_at` IS NULL ORDER BY `id` DESC
	//err := gonet.Gorm().Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).Find(&users).Error

	var total int64
	offset := (page - 1) * size

	// SELECT count(*) FROM `user` WHERE `deleted_at` IS NULL
	err := gonet.Gorm().Client().Model(u).Count(&total).Error
	if err != nil || total == 0 {
		return nil, 0, err
	}
	// SELECT `id`,`uname`,`email`,`created_at` FROM `user` WHERE `deleted_at` IS NULL ORDER BY `id` DESC LIMIT ?
	err = gonet.Gorm().Client().
		Select([]string{"id", "uname", "email", "created_at"}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Offset(offset).
		Limit(size).
		Find(&users).Error
	return users, total, err
}

func (u *User) Update(user *User) error {
	// UPDATE `user` SET `uname`=?,`updated_at`=? WHERE `id` = ? AND `deleted_at` IS NULL
	return gonet.Gorm().Client().
		Model(&User{}).
		Select([]string{"uname"}).
		Where("id", user.ID).
		Updates(user).Error
}
