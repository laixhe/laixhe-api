package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// UserRouter 用户相关
func UserRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("user")
	c := controllers.NewUser()
	// not token
	//{
	//}

	// token auto
	//jwtAutoRouter := groupRouter.Use(xgin.JwtAuthAuto(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	//{
	//}

	// token
	jwtRouter := groupRouter.Use(xgin.JwtAuth(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	{
		jwtRouter.GET("info", c.Info)      // 用户信息
		jwtRouter.GET("list", c.List)      // 用户列表
		jwtRouter.POST("update", c.Update) // 修改用户信息
	}
}
