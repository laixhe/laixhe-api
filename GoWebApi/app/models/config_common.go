package models

import (
	"gorm.io/gorm"
)

// ConfigCommonTable 通用配置表名
const ConfigCommonTable = "config_common"

// ConfigCommon 通用配置
type ConfigCommon struct {
	ID       int    `gorm:"column:id;type:int;autoIncrement;primaryKey"`
	Key      string `gorm:"column:key;type:string;size:255;not null;default:'';index"`
	Value    string `gorm:"column:value;type:string;size:512;not null;default:''"`
	Describe string `gorm:"column:describe;type:string;size:255;not null;default:'';comment:描述"`
}

func (m *ConfigCommon) TableName() string {
	return ConfigCommonTable
}

// List 查询通用配置列表，可选按 key 过滤
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
