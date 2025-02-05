package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"

	"webapi/app/controllers"
	"webapi/core"
)

// UserRouter 用户相关
func UserRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("user")
	cs := controllers.NewUser()
	// not token
	//{
	//}

	// token auto
	//jwtAutoRouter := groupRouter.Group("", xgin.JwtAuthAuto(core.Config().Jwt, core.ErrorAuthInvalid(nil)))
	//{
	//}

	// token
	jwtRouter := groupRouter.Group("", xgin.JwtAuth(core.Config().Jwt, &xerror.Error{}))
	{
		jwtRouter.GET("info", cs.Info)      // 用户信息
		jwtRouter.GET("list", cs.List)      // 用户列表
		jwtRouter.POST("update", cs.Update) // 修改用户信息
	}
}
