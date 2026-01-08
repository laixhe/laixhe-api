package models

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
