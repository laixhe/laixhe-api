package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// AuthRouter 鉴权相关
func AuthRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("auth")
	cs := controllers.NewAuth()

	// not token
	{
		groupRouter.POST("register", cs.Register) // 注册
		groupRouter.POST("login", cs.Login)       // 登录
	}

	// token auto
	//jwtAutoRouter := groupRouter.Group("", xgin.JwtAuthAuto(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	//{
	//}

	// token
	jwtRouter := groupRouter.Group("", xgin.JwtAuth(core.Config().Jwt, &xerror.Error{}))
	{
		jwtRouter.POST("refresh", cs.Refresh) // 刷新Jwt
	}
}
