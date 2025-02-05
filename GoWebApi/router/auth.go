package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// AuthRouter 鉴权相关
func AuthRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("auth")
	c := controllers.NewAuth()
	// not token
	{
		groupRouter.POST("register", c.Register) // 注册
		groupRouter.POST("login", c.Login)       // 登录
	}

	// token auto
	//jwtAutoRouter := groupRouter.Group("", xgin.JwtAuthAuto(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	//{
	//}

	// token
	jwtRouter := groupRouter.Group("", xgin.JwtAuth(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	{
		jwtRouter.POST("refresh", c.Refresh) // 刷新Jwt
	}
}
