package core

import (
	"errors"
	"fmt"

	"github.com/laixhe/gonet/config"
	"github.com/laixhe/gonet/db/gorm/orm"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/xlog"
)

// Addr HTTP 服务监听地址
type Addr struct {
	IP      string `mapstructure:"ip"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"` // 请求超时时间(单位秒), 缺省 30 秒, 用于请求超时中间件
}

// Address 返回 "ip:port" 格式的监听地址
func (a *Addr) Address() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

// Common 运行时通用配置，从数据库动态加载
type Common struct {
	Env string
}

// Limit 接口限流配置
type Limit struct {
	Enable bool `mapstructure:"enable"`
	Max    int  `mapstructure:"max"`
	Window int  `mapstructure:"window"`
}

// Config 配置
type Config struct {
	Http   *Addr        `mapstructure:"http"`
	Log    *xlog.Config `mapstructure:"log"`
	Orm    *orm.Config  `mapstructure:"orm"`
	Jwt    *jwt.Config  `mapstructure:"jwt"`
	Limit  *Limit       `mapstructure:"limit"`
	Common *Common      `mapstructure:"-"`
}

// Check 校验配置有效性，缺省日志配置自动补全
func (c *Config) Check() error {
	if c.Http == nil {
		return errors.New("http config is nil")
	}
	if c.Http.Port <= 0 {
		return errors.New("http port is invalid")
	}
	// 请求超时缺省 30 秒, 用于请求超时中间件;
	// config.yaml 已显式写 timeout: 30, 此缺省兜底仅在配置缺失/非法时生效
	if c.Http.Timeout <= 0 {
		c.Http.Timeout = 30
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
	if c.Limit == nil {
		// 缺省配置: 启用, 单个 IP 在 60 秒内最多 1000 次
		c.Limit = &Limit{Enable: true, Max: 1000, Window: 60}
	}
	return nil
}

// NewConfig 加载配置文件并校验
func NewConfig(configFile string) *Config {
	c := &Config{
		Common: &Common{},
	}
	config.Init(configFile, c)
	if err := c.Check(); err != nil {
		panic(err)
	}
	return c
}
