package services

import (
	"context"

	"webapi/app/models"
	"webapi/core"
)

// Service 业务服务逻辑
type Service struct {
	server *core.Server
	Auth   *Auth
	User   *User
}

// NewService 创建业务服务实例，初始化子模块并加载运行时配置
func NewService(server *core.Server) *Service {
	service := &Service{
		server: server,
		Auth:   NewAuth(server),
		User:   NewUser(server),
	}
	service.initConfigCommon()
	return service
}

// initConfigCommon 从数据库 config_common 表加载运行时配置（如环境标识 env）
func (s *Service) initConfigCommon() {
	configs, err := new(models.ConfigCommon).List(s.server.Gorm(context.Background()))
	if err != nil {
		s.server.Log().Errorf("initConfigCommon failed: %v", err)
		return
	}
	for _, v := range configs {
		if v.Key == models.ConfigCommonEnv {
			s.server.Config().Common.Env = v.Value
		}
	}

	s.server.Log().Debugf("config http=%#v", s.server.Config().Http)
	s.server.Log().Debugf("config common=%#v", s.server.Config().Common)
}
