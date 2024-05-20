package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"webapi/app/result"
	"webapi/core/config"
	"webapi/core/errorx"
	"webapi/core/jwtx"
	"webapi/core/logx"
	pbCode "webapi/profile/gen/code"
)

// Authorization
const (
	Authorization = "Authorization"
	Bearer        = "Bearer "
	BearerLen     = 7
)

// JwtAuth 鉴权
func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ecode := pbCode.ECode_AuthNotLogin
		token := ctx.Request.Header.Get(Authorization)
		if len(token) > 0 {
			if strings.HasPrefix(token, Bearer) {
				claims, err := jwtx.ParseToken(token[BearerLen:], config.Get().Jwt)
				if err == nil {
					ctx.Set(jwtx.AuthUid, claims.Uid)
					ctx.Next()
					return
				} else {
					if errors.Is(err, jwtx.ErrTokenExpired) {
						ecode = pbCode.ECode_AuthExpire
					}
					logx.Errorf("authorization:%v error:%v", token, err)
				}
			}
		}
		result.ResponseError(ctx, errorx.NewError(ecode, nil))
		// 返回错误
		ctx.Abort()
	}
}

func Uid(c *gin.Context) (uint64, *errorx.Error) {
	value, exists := c.Get(jwtx.AuthUid)
	if exists {
		return value.(uint64), nil
	}
	return 0, errorx.NewError(pbCode.ECode_AuthNotLogin, nil)
}
