package core

import (
	"errors"
	"fmt"
	"os"

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

// Bcrypt bcrypt 计算 worker 池配置
type Bcrypt struct {
	// Workers worker goroutine 数量 (并发执行的密码哈希/校验数)
	// 0 或负数时自动取 GOMAXPROCS (CPU 逻辑核数); bcrypt 为纯 CPU 计算, 建议不超过核数
	Workers int `mapstructure:"workers"`
}

// Config 配置
type Config struct {
	Http   *Addr        `mapstructure:"http"`
	Log    *xlog.Config `mapstructure:"log"`
	Orm    *orm.Config  `mapstructure:"orm"`
	Jwt    *jwt.Config  `mapstructure:"jwt"`
	Limit  *Limit       `mapstructure:"limit"`
	Bcrypt *Bcrypt      `mapstructure:"bcrypt"`
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
	if c.Bcrypt == nil {
		// 缺省配置: workers=0 自动取 GOMAXPROCS (bcrypt 为纯 CPU 计算, 池大于核数无吞吐收益)
		c.Bcrypt = &Bcrypt{Workers: 0}
	}
	return nil
}

// NewConfig 加载配置文件并校验
//
// JWT 密钥支持环境变量 JWT_SECRET_KEY 覆盖 config.yaml,
// 生产环境应通过环境变量注入, 避免把密钥写死在配置仓库里。
func NewConfig(configFile string) *Config {
	c := &Config{
		Common: &Common{},
	}
	// 配置文件缺失/损坏时直接以明确文案启动失败, 避免后续以晦涩的 panic("http config is nil") 崩溃
	if err := config.Init(configFile, c); err != nil {
		panic(fmt.Errorf("load config file %q: %w", configFile, err))
	}
	// 环境变量优先于配置文件 (需在 Check 之前完成覆盖);
	// 先判 nil: 防止 config.yaml 缺失 jwt 段时 c.Jwt 为空直接空指针 panic
	if secret := os.Getenv("JWT_SECRET_KEY"); secret != "" {
		if c.Jwt == nil {
			panic("jwt config is nil (config.yaml 缺少 jwt 段)")
		}
		c.Jwt.SecretKey = secret
	}
	if err := c.Check(); err != nil {
		panic(err)
	}
	return c
}
