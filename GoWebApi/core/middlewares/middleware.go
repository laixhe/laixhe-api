package middlewares

import (
	"time"

	contribJwt "github.com/gofiber/contrib/v3/jwt"

	"webapi/core"
)

// Middleware 中间件: JWT 鉴权与 IP 限流
type Middleware struct {
	// server 持有服务实例, 供限流触发时输出告警日志
	server *core.Server
	// UseJwtConfig 强制 JWT 校验，无 Token 返回 401
	UseJwtConfig contribJwt.Config
	// rateLimiter IP 限流器 (基于客户端 IP 的滑动窗口)
	rateLimiter *RateLimiter
	// limitConfig 接口限流配置
	limitConfig *core.Limit
}

// NewMiddleware 创建中间件，包含强制 JWT 校验与 IP 限流器
func NewMiddleware(server *core.Server) *Middleware {
	return &Middleware{
		server: server,
		UseJwtConfig: contribJwt.Config{
			SigningKey: contribJwt.SigningKey{Key: []byte(server.Config().Jwt.SecretKey)},
			Claims:     &JwtClaims{},
		},
		rateLimiter: NewRateLimiter(server.Config().Limit.Max, time.Duration(server.Config().Limit.Window)*time.Second),
		limitConfig: server.Config().Limit,
	}
}
