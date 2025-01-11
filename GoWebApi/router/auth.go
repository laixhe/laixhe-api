package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// AuthRouter 鉴权相关
func AuthRouter(r *gin.RouterGroup) {
	authRouter := r.Group("auth")
	c := controllers.NewAuth()
	// not token
	{
		authRouter.POST("register", c.Register) // 注册
		authRouter.POST("login", c.Login)       // 登录
	}
	// token
	authRouterJwt := authRouter.Use(xgin.JwtAuth(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	{
		authRouterJwt.POST("refresh", c.Refresh) // 刷新Jwt
	}
}
