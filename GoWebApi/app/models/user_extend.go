package models

const UserExtendTable = "user_extend"

// UserExtend 用户扩展
type UserExtend struct {
	Uid           int    `gorm:"column:uid;type:int;not null;primaryKey;comment:用户UID"`
	WechatUnionid string `gorm:"column:wechat_unionid;type:string;size:255;not null;default:'';comment:微信unionid"`
	WechatOpenid  string `gorm:"column:wechat_openid;type:string;size:255;not null;default:'';index;comment:微信openid"`
}

func (*UserExtend) TableName() string {
	return UserExtendTable
}
