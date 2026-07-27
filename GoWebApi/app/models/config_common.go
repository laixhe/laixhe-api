package models

import (
	"gorm.io/gorm"
)

const ConfigCommonTable = "config_common"

// ConfigCommon 通用配置
type ConfigCommon struct {
	ID       int    `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	Key      string `gorm:"column:key;type:string;size:50;not null;default:''"`
	Value    string `gorm:"column:value;type:string;size:500;not null;default:''"`
	Describe string `gorm:"column:describe;type:string;size:500;not null;default:'';comment:描述"`
}

func (m *ConfigCommon) TableName() string {
	return ConfigCommonTable
}

func (m *ConfigCommon) List(db *gorm.DB, keys ...string) ([]ConfigCommon, error) {
	var list []ConfigCommon
	if len(keys) > 0 {
		db = db.Where("key IN ?", keys)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
