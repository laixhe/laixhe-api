package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// UserRouter 用户相关
func UserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("user")
	c := controllers.NewUser()
	// not token
	// token
	userRouterJwt := userRouter.Use(xgin.JwtAuth(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	{
		userRouterJwt.GET("info", c.Info)      // 用户信息
		userRouterJwt.GET("list", c.List)      // 用户列表
		userRouterJwt.POST("update", c.Update) // 修改用户信息
	}
}
