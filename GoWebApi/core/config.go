package core

import (
	"errors"
	"fmt"

	"github.com/laixhe/gonet/config"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/orm/orm"
	"github.com/laixhe/gonet/xlog"
)

type Addr struct {
	IP      string `mapstructure:"ip"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

func (a *Addr) Addr() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

type Common struct {
	Env string
}

// Config 配置
type Config struct {
	Http   *Addr        `mapstructure:"http"`
	Log    *xlog.Config `mapstructure:"log"`
	Orm    *orm.Config  `mapstructure:"orm"`
	Jwt    *jwt.Config  `mapstructure:"jwt"`
	Common *Common      `mapstructure:"-"`
}

func (c *Config) Check() error {
	if c.Http == nil {
		return errors.New("http config is nil")
	}
	if c.Http.Port <= 0 {
		return errors.New("http port is invalid")
	}
	if c.Log == nil {
		c.Log = &xlog.Config{
			Run: xlog.RunTypeConsole,
		}
	}
	if err := c.Orm.Check(); err != nil {
		return err
	}
	if err := c.Jwt.Check(); err != nil {
		return err
	}
	return nil
}

func NewConfig(configFile string) *Config {
	c := &Config{}
	config.Init(configFile, c)
	if err := c.Check(); err != nil {
		panic(err)
	}
	return c
}
