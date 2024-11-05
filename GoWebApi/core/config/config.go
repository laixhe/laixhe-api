package config

import (
	"github.com/laixhe/gonet/configx"
	"github.com/laixhe/gonet/jwtx"
	"github.com/laixhe/gonet/logx"
	"github.com/laixhe/gonet/proto/gen/config/capp"
	"github.com/laixhe/gonet/proto/gen/config/cauth"
	"github.com/laixhe/gonet/proto/gen/config/cgorm"
	"github.com/laixhe/gonet/proto/gen/config/clog"
	"github.com/laixhe/gonet/proto/gen/config/cmongodb"
	"github.com/laixhe/gonet/proto/gen/config/credis"
	"github.com/laixhe/gonet/proto/gen/config/cserver"
)

type Config struct {
	App     *capp.App         `mapstructure:"app"`
	Http    *cserver.Server   `mapstructure:"http"`
	Log     *clog.Log         `mapstructure:"log"`
	Gorm    *cgorm.Gorm       `mapstructure:"gorm"`
	Mongodb *cmongodb.MongoDB `mapstructure:"mongodb"`
	Redis   *credis.Redis     `mapstructure:"redis"`
	Jwt     *cauth.Jwt        `mapstructure:"jwt"`
}

func Init(configFile string) *Config {
	c := &Config{}
	configx.Init(configFile, false, c)
	logx.Init(c.Log)
	return c
}

// AppChecking 检查App配置
func (c *Config) AppChecking() *Config {
	if c.App == nil {
		panic("app config is nil")
	}
	if c.App.Version == "" {
		c.App.Version = "v0.1"
	}
	if c.App.Mode == "" {
		c.App.Mode = capp.ModeType_debug.String()
	} else {
		c.App.Mode = capp.ModeType_name[capp.ModeType_value[c.App.Mode]]
	}
	logx.Debugf("app config=%v", c.App)
	return c
}

// HttpChecking 检查Http配置
func (c *Config) HttpChecking() *Config {
	if c.Http == nil {
		panic("http config is nil")
	}
	if c.Http.Port <= 0 || c.Http.Port > 65535 {
		panic("http config port error: 1~65535")
	}
	logx.Debugf("http config=%v", c.Http)
	return c
}

// JwtChecking 检查Jwt配置
func (c *Config) JwtChecking() *Config {
	if err := jwtx.Checking(c.Jwt); err != nil {
		panic(err)
	}
	logx.Debugf("jwt config=%v", c.Jwt)
	return c
}
