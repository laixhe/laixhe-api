package core

import (
	"github.com/laixhe/gonet/orm/mysql"
	"github.com/laixhe/gonet/orm/orm"
	"github.com/laixhe/gonet/xfiber"
	"github.com/laixhe/gonet/xlog"
)

// DEFAULT 默认key
const DEFAULT = "default"

// Server 服务
type Server struct {
	config *Config
	log    *xlog.LogClient
	server *xfiber.Server
	orm    map[string]orm.Client
}

// NewServer 创建服务
func NewServer(configFile string) *Server {
	config := NewConfig(configFile)
	// 初始化日志
	config.Log.CallerSkip = 1
	logClient, err := xlog.Init(config.Log)
	if err != nil {
		panic(err)
	}
	server := xfiber.New(logClient.Logger()).
		UseCors().
		UseRecover()
	s := &Server{
		config: config,
		log:    logClient,
		server: server,
		orm:    make(map[string]orm.Client),
	}
	return s.init()
}

func (s *Server) Server() *xfiber.Server {
	return s.server
}

func (s *Server) Config() *Config {
	return s.config
}

func (s *Server) Log() *xlog.LogClient {
	return s.log
}

func (s *Server) initOrm(config *orm.Config, key ...string) error {
	db, err := mysql.Init(config, NewOrmWriter(s.server.LoggerConfig()), xfiber.RequestIdLogKey)
	if err != nil {
		return err
	}
	if len(key) > 0 {
		s.orm[key[0]] = db
	} else {
		s.orm[DEFAULT] = db
	}
	return nil
}

func (s *Server) Orm(key ...string) orm.Client {
	if len(key) > 0 {
		return s.orm[key[0]]
	}
	return s.orm[DEFAULT]
}

func (s *Server) init() *Server {
	if err := s.initOrm(s.config.Orm); err != nil {
		panic(err)
	}
	return s
}
