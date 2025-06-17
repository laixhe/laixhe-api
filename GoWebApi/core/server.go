package core

import (
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/orm"
	"github.com/laixhe/gonet/xlog"
	"gorm.io/gorm/logger"
)

// DEFAULT 默认key
const DEFAULT = "default"

// Server 服务器
type Server struct {
	config *Config
	app    *fiber.App
	log    *fiberzap.LoggerConfig
	orm    map[string]*orm.GormClient
}

func NewServer(configFile string) *Server {
	// 初始化配置
	config := NewConfig(configFile)
	if config.RequestIdKey == "" {
		config.RequestIdKey = "requestId"
	}
	// 初始化日志
	config.Log.CallerSkip = 2
	zapClient, err := xlog.Init(config.Log)
	if err != nil {
		panic(err)
	}
	zapLogger := fiberzap.NewLogger(fiberzap.LoggerConfig{
		ExtraKeys: []string{config.RequestIdKey},
		SetLogger: zapClient.Logger(),
	})
	log.SetLogger(zapLogger) // 替换默认日志
	return &Server{
		config: config,
		log:    zapLogger,
		orm:    make(map[string]*orm.GormClient),
	}
}

func (s *Server) Config() *Config {
	return s.config
}

func (s *Server) Log() *fiberzap.LoggerConfig {
	return s.log
}

func (s *Server) initOrm(config *orm.Config, key ...string) error {
	logLevel := logger.Info
	if config.LogLevel == logger.Silent ||
		config.LogLevel == logger.Error ||
		config.LogLevel == logger.Warn ||
		config.LogLevel == logger.Info {
		logLevel = config.LogLevel
	}
	dbLogger := orm.NewOrmLogger(NewOrmWriter(s.log), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logLevel,
	}, s.config.RequestIdKey)
	db, err := orm.Init(config, dbLogger)
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

func (s *Server) Orm(key ...string) *orm.GormClient {
	if len(key) > 0 {
		return s.orm[key[0]]
	} else {
		return s.orm[DEFAULT]
	}
}

func (s *Server) Init(app *fiber.App) {
	if err := s.initOrm(s.config.Orm); err != nil {
		panic(err)
	}
	if err := jwt.CheckConfig(s.config.Jwt); err != nil {
		panic(err)
	}
	s.app = app
}

// Listen 启动服务
func (s *Server) Listen() error {
	return s.app.Listen(s.config.Addr())
}
