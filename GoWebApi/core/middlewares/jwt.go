package middlewares

import (
	"context"
	"errors"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const JwtKey = "jwt"

type JwtContextUidKey struct{}

// JwtClaims 可根据业务自行添加字段
type JwtClaims struct {
	Uid int `json:"uid"`
	jwtv5.RegisteredClaims
}

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

func (claims *JwtClaims) GetUid() int {
	return claims.Uid
}

// UseJwt JWT中间件
func UseJwt(secretKey string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secretKey)},
		ContextKey:   JwtKey,
		Claims:       &JwtClaims{},
		ErrorHandler: UseJwtErrorHandler,
	})
}

// UseJwtErrorHandler 自定义JWT错误响应
func UseJwtErrorHandler(ctx *fiber.Ctx, err error) error {
	authorization := ctx.Get(fiber.HeaderAuthorization)
	log.WithContext(ctx.UserContext()).Errorf("jwt: %s error: %v", authorization, err)
	return ctx.Next()
	// if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
	// }
	// return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.NewError(fiber.StatusUnauthorized, err.Error()))
}

// UseJwtClaims 获取JWT中的claims
func UseJwtClaims(isSkipError ...bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims, err := GetJwtClaims(ctx)
		if err != nil {
			if len(isSkipError) > 0 && isSkipError[0] {
				return ctx.Next()
			}
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.NewError(fiber.StatusUnauthorized, err.Error()))
		}
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), JwtContextUidKey{}, claims.Uid))
		return ctx.Next()
	}
}

// GetJwtClaims 获取JWT中的claims
func GetJwtClaims(ctx *fiber.Ctx) (*JwtClaims, error) {
	token, isToken := ctx.Locals(JwtKey).(*jwtv5.Token)
	if isToken {
		claims, isClaims := token.Claims.(*JwtClaims)
		if isClaims {
			if claims.Uid > 0 {
				return claims, nil
			}
		}
	}
	return nil, errors.New("尚未登录")
}

// GetUid 获取当前请求的用户ID
func GetUid(ctx *fiber.Ctx) (int, error) {
	uid, is := ctx.UserContext().Value(JwtContextUidKey{}).(int)
	if is && uid > 0 {
		return uid, nil
	}
	return 0, fiber.NewError(fiber.StatusUnauthorized, "尚未登录")
}
