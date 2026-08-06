package middlewares

import (
	contribJwt "github.com/gofiber/contrib/v3/jwt"
	"github.com/laixhe/gonet/xfiber"
)

// Middleware JWT 中间件配置
type Middleware struct {
	// UseJwtConfig 强制 JWT 校验，无 Token 返回 401
	UseJwtConfig contribJwt.Config
	// UseJwtConfigNext 可选 JWT 校验，无 Token 时继续传递请求（用于可选认证场景）
	UseJwtConfigNext contribJwt.Config
}

// NewMiddleware 创建 JWT 中间件配置，包含强制校验和可选校验两种模式
func NewMiddleware(jwtSecretKey string) *Middleware {
	return &Middleware{
		UseJwtConfig: contribJwt.Config{
			SigningKey: contribJwt.SigningKey{Key: []byte(jwtSecretKey)},
			Claims:     &JwtClaims{},
		},
		UseJwtConfigNext: contribJwt.Config{
			ErrorHandler: xfiber.JwtErrorHandlerNext,
			SigningKey:   contribJwt.SigningKey{Key: []byte(jwtSecretKey)},
			Claims:       &JwtClaims{},
		},
	}
}
