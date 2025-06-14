package core

import (
	"strconv"

	"github.com/laixhe/gonet/config"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/orm"
	"github.com/laixhe/gonet/xlog"
)

// HttpConfig http 配置
type HttpConfig struct {
	IP      string `mapstructure:"ip"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

// Config 配置
type Config struct {
	Http         *HttpConfig  `mapstructure:"http"`
	Log          *xlog.Config `mapstructure:"log"`
	Orm          *orm.Config  `mapstructure:"orm"`
	Jwt          *jwt.Config  `mapstructure:"jwt"`
	RequestIdKey string       `mapstructure:"request_id_key"`
}

func NewConfig(configFile string) *Config {
	data := &Config{}
	config.Init(configFile, data)
	return data
}

// Addr 获取 http 地址
func (c *Config) Addr() string {
	return c.Http.IP + ":" + strconv.Itoa(c.Http.Port)
}
