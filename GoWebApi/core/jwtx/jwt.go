package jwtx

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

/**

jwt:
  secret_key: 6Kbj0VFeXYMp60lEyiFoVq4UzqX8Z0GSSfnvTh2VuAQn0oHgQNYexU6yYVTk4xf9
  # 过期时长(单位秒)
  expire_time: 604800

*/

var (
	ErrTokenExpired = errors.New("token is expired") // 令牌已过期
	ErrTokenInvalid = errors.New("token invalid")    // 令牌无效
)

const (
	AuthUid = "uid" // CustomClaims.Uid
)

// Config jwt
type Config struct {
	SecretKey  string `mapstructure:"secret_key"`  // jwt secret key
	ExpireTime int    `mapstructure:"expire_time"` // 过期时长(单位秒)
}

// CustomClaims 自定义声明类型 并内嵌 jwt.RegisteredClaims
// jwt 包自带的 jwt.RegisteredClaims 只包含了官方字段
type CustomClaims struct {
	// 可根据需要自行添加字段
	Uid uint64 `json:"uid"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(uid uint64, c *Config) (string, error) {
	claims := CustomClaims{
		Uid: uid,
	}

	nowTime := time.Now()
	// 过期时间
	if c.ExpireTime > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(nowTime.Add(time.Duration(c.ExpireTime) * time.Second))
	}
	// 发布时间（创建时间）
	claims.IssuedAt = jwt.NewNumericDate(nowTime)
	// 生效时间
	claims.NotBefore = jwt.NewNumericDate(nowTime)

	// 使用指定的签名方法创建签名对象（使用 HS256 签名算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString([]byte(c.SecretKey))
}

// ParseToken 解析JWT
func ParseToken(tokenString string, c *Config) (*CustomClaims, error) {
	// 如果是自定义 Claim 结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(c.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrTokenInvalid
		}
		return nil, ErrTokenInvalid
	}
	if token != nil {
		// 对 token 对象中的 Claim 进行类型断言
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验 token
			return claims, nil
		}
	}
	return nil, ErrTokenInvalid
}
