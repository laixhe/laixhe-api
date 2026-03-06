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
