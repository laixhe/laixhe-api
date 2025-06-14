package models

import (
	"webapi/core"

	"github.com/laixhe/gonet/orm"
)

type Model struct {
	server     *core.Server
	User       *User
	UserExtend *UserExtend
}

func NewModel(server *core.Server) *Model {
	return &Model{
		server:     server,
		User:       &User{},
		UserExtend: &UserExtend{},
	}
}

func (m *Model) orm() *orm.GormClient {
	return m.server.Orm()
}
