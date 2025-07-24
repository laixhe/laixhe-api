package core

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	golog "github.com/laixhe/gonet/log"
	"github.com/laixhe/gonet/orm/mysql"
	"github.com/laixhe/gonet/orm/orm"

	"webapi/core/middlewares"
)

// DEFAULT 默认key
const DEFAULT = "default"

// RequestIdKey 请求ID key
const RequestIdKey = "requestId"

// Server 服务
type Server struct {
	config *Config
	log    *fiberzap.LoggerConfig
	app    *fiber.App
	orm    map[string]orm.Client
}

// NewServer 创建服务
func NewServer(configFile string) *Server {
	config := NewConfig(configFile)
	// 初始化日志
	config.Log.CallerSkip = 2
	zapClient, err := golog.Init(config.Log)
	if err != nil {
		panic(err)
	}
	zapLogger := fiberzap.NewLogger(fiberzap.LoggerConfig{
		ExtraKeys: []string{RequestIdKey},
		SetLogger: zapClient.Logger(),
	})
	// 替换默认日志
	log.SetLogger(zapLogger)
	// http服务
	app := fiber.New(fiber.Config{
		// 默认错误处理
		ErrorHandler: middlewares.UseErrorDefault(),
	})
	s := &Server{
		config: config,
		log:    zapLogger,
		app:    app,
		orm:    make(map[string]orm.Client),
	}
	return s.init()
}

// Config 获取配置
func (s *Server) Config() *Config {
	return s.config
}

// Log 获取日志
func (s *Server) Log() *fiberzap.LoggerConfig {
	return s.log
}

// 初始化ORM
func (s *Server) initOrm(config *orm.Config, key ...string) error {
	db, err := mysql.Init(config, NewOrmWriter(s.log), RequestIdKey)
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

// Orm 获取ORM
func (s *Server) Orm(key ...string) orm.Client {
	if len(key) > 0 {
		return s.orm[key[0]]
	} else {
		return s.orm[DEFAULT]
	}
}

// HttpMiddleware 中间件
func (s *Server) HttpMiddleware(fn func(app *fiber.App)) {
	fn(s.app)
}

// HttpGroup 路由分组
func (s *Server) HttpGroup(prefix string) fiber.Router {
	return s.app.Group(prefix)
}

// HttpStart 启动Http服务
func (s *Server) HttpStart() error {
	return s.app.Listen(s.config.Http.Addr())
}

// Init 初始化
func (s *Server) init() *Server {
	if err := s.initOrm(s.config.Orm); err != nil {
		panic(err)
	}
	return s
}
