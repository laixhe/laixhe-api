package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"

	"webapi/api/gen/ecode"
	"webapi/app/controllers"
	"webapi/core"
)

// AuthRouter 鉴权相关
func AuthRouter(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	c := controllers.NewAuth()

	// not token

	auth.POST("/register", c.Register) // 注册
	auth.POST("/login", c.Login)       // 登陆

	// token

	jwt := auth.Use(xgin.JwtAuth(core.Config().Jwt, core.NewError(ecode.ECode_AuthInvalid, nil)))
	jwt.POST("/refresh", c.Refresh) // 刷新Jwt
}
