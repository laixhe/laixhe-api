package middlewares

import (
	"time"

	contribJwt "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// JwtClaims JWT 令牌载荷，存储用户 UID
type JwtClaims struct {
	Uid int `json:"uid"`
	jwtv5.RegisteredClaims
}

// GetUid 返回 JWT 中的用户 UID
func (jc *JwtClaims) GetUid() int {
	return jc.Uid
}

// NewJwtClaims 创建 JWT 载荷并设置过期时间、发布时间、生效时间
func NewJwtClaims(uid int, expireTime int64) *JwtClaims {
	custom := &JwtClaims{
		Uid: uid,
	}
	nowTime := time.Now()
	custom.ExpiresAt = jwtv5.NewNumericDate(nowTime.Add(time.Duration(expireTime) * time.Second)) // 过期时间
	custom.IssuedAt = jwtv5.NewNumericDate(nowTime)                                               // 发布时间（创建时间）
	custom.NotBefore = jwtv5.NewNumericDate(nowTime)                                              // 生效时间
	return custom
}

// GetJwtClaims 从 Fiber 上下文中提取并验证 JWT 载荷
func GetJwtClaims(ctx fiber.Ctx) (*JwtClaims, error) {
	token := contribJwt.FromContext(ctx)
	if token != nil {
		claims, isClaims := token.Claims.(*JwtClaims)
		if isClaims {
			if claims.Uid > 0 {
				return claims, nil
			}
		}
	}
	return nil, fiber.NewError(fiber.StatusUnauthorized)
}
