package core

import (
	"context"

	"github.com/laixhe/gonet/db/gorm/mysql"
	"github.com/laixhe/gonet/db/gorm/orm"
	"github.com/laixhe/gonet/xfiber"
	"github.com/laixhe/gonet/xlog"
	"gorm.io/gorm"
)

// DEFAULT 默认key
const DEFAULT = "default"

// Server 服务
type Server struct {
	config *Config
	log    *xlog.ZClient
	server *xfiber.Server
	orm    map[string]orm.Client
}

// NewServer 创建服务
func NewServer(configFile string) *Server {
	config := NewConfig(configFile)
	// 初始化日志
	config.Log.CallerSkip = 1
	logClient, err := xlog.InitZap(config.Log)
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

// Server 返回 Fiber 服务器实例
func (s *Server) Server() *xfiber.Server {
	return s.server
}

// Config 返回应用配置
func (s *Server) Config() *Config {
	return s.config
}

// Log 返回日志客户端
func (s *Server) Log() *xlog.ZClient {
	return s.log
}

// initOrm 初始化 ORM 数据库连接，支持环境变量展开 DSN
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

// Orm 返回 ORM 客户端，可选指定 key（默认 "default"）
func (s *Server) Orm(key ...string) orm.Client {
	if len(key) > 0 {
		return s.orm[key[0]]
	}
	return s.orm[DEFAULT]
}

// Gorm 返回绑定了 context 的 GORM 实例
func (s *Server) Gorm(ctx context.Context, key ...string) *gorm.DB {
	return s.Orm(key...).WithContext(ctx)
}

// init 初始化服务（目前仅初始化 ORM）
func (s *Server) init() *Server {
	if err := s.initOrm(s.config.Orm); err != nil {
		panic(err)
	}
	return s
}
