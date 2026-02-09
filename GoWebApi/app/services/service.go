package services

import (
	"context"
	"fmt"

	"webapi/app/models"
	"webapi/core"
)

// Service 业务服务逻辑
type Service struct {
	server *core.Server
	Auth   *Auth
	User   *User
}

func NewService(server *core.Server) *Service {
	service := &Service{
		server: server,
		Auth:   NewAuth(server),
		User:   NewUser(server),
	}
	service.initConfigCommon()
	return service
}

func (s *Service) initConfigCommon() {
	configs, err := new(models.ConfigCommon).List(s.server.Gorm(context.Background()))
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
