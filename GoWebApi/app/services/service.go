package services

import (
	"context"
	"fmt"

	"webapi/app/models"
	"webapi/app/models/dao"
	"webapi/core"
)

// Service 业务服务逻辑
type Service struct {
	server *core.Server
	Auth   *Auth
	User   *User
}

func NewService(server *core.Server, modelDao *dao.Dao) *Service {
	service := &Service{
		server: server,
		Auth:   NewAuth(server, modelDao),
		User:   NewUser(server, modelDao),
	}
	return service
}

func (s *Service) initConfigCommon(modelDao *dao.Dao) {
	configs, err := modelDao.ListConfigCommon(context.Background())
	if err != nil {
		panic(err)
	}
	for _, v := range configs {
		if v.Key == models.ConfigCommonEnv {
			s.server.Config().Common.Env = v.Value
		}
	}

	fmt.Printf("config http=%#v\n", s.server.Config().Http)
	fmt.Printf("config log=%#v\n", s.server.Config().Log)
	fmt.Printf("config orm=%#v\n", s.server.Config().Orm)
	fmt.Printf("config jwt=%#v\n", s.server.Config().Jwt)
	fmt.Printf("config common=%#v\n", s.server.Config().Common)
}
