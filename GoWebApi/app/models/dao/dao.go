package dao

import (
	"context"

	"github.com/laixhe/gonet/orm/orm"
	"gorm.io/gorm"

	"webapi/app/models"
	"webapi/core"
)

// Dao 业务数据操作
type Dao struct {
	server         *core.Server
	ConfigCommon   *models.ConfigCommon
	User           *models.User
	UserExtend     *models.UserExtend
	UserThirdParty *models.UserThirdParty
}

func NewDao(server *core.Server) *Dao {
	return &Dao{
		server:         server,
		ConfigCommon:   &models.ConfigCommon{},
		User:           &models.User{},
		UserExtend:     &models.UserExtend{},
		UserThirdParty: &models.UserThirdParty{},
	}
}

func (d *Dao) Orm() orm.Client {
	return d.server.Orm()
}

func (d *Dao) WithContext(ctx context.Context) *gorm.DB {
	return d.server.Orm().WithContext(ctx)
}

// GetById 以 id 获取数据
// data 指针传递的结构(表结构)
func (d *Dao) GetById(ctx context.Context, id int, data any) error {
	return d.WithContext(ctx).Where("id", id).Take(data).Error
}

// GetByKeyValue 以 id 获取数据
// key   要查询的字段名
// value 要查询的字段名的值
// data 指针传递的结构(表结构)
func (d *Dao) GetByKeyValue(ctx context.Context, key string, value, data any) error {
	return d.WithContext(ctx).Where(key, value).Take(data).Error
}

// GetByUid 以 uid 获取数据
// data 指针传递的结构(表结构)
func (d *Dao) GetByUid(ctx context.Context, uid int, data any) error {
	return d.WithContext(ctx).Where("uid", uid).First(data).Error
}

// Save 会保存所有的字段，即使字段是零值
// data 指针传递的结构(表结构)
func (d *Dao) Save(ctx context.Context, data any) error {
	return d.WithContext(ctx).Save(data).Error
}

// Create 创建数据
// data 指针传递的结构或者数组结构(表结构)
func (d *Dao) Create(ctx context.Context, data any) error {
	return d.WithContext(ctx).Create(data).Error
}

// Delete 删除数据
// data 指针传递的结构或者数组结构(表结构)（必须包含 id 字段并赋值）
func (d *Dao) Delete(ctx context.Context, data any) error {
	return d.WithContext(ctx).Delete(data).Error
}

// DeleteById 以 id 删除数据
// model 指针传递的结构(表结构)
func (d *Dao) DeleteById(ctx context.Context, id int, model any) error {
	return d.WithContext(ctx).Where("id", id).Delete(model).Error
}

// UpdatesById 以 id 修改数据
// model 指针传递的结构(表结构)
// data  修改的数据(表对应的字段)
func (d *Dao) UpdatesById(ctx context.Context, id int, model any, data map[string]any) error {
	return d.WithContext(ctx).Model(model).Where("id", id).Updates(data).Error
}
