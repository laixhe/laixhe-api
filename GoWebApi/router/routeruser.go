package router

import (
	"github.com/gin-gonic/gin"

	"webapi/app/controllers"
	"webapi/core/ginx"
)

// UserRouter 用户相关
func UserRouter(r *gin.RouterGroup) {
	user := r.Group("user")
	c := controllers.NewUser()

	// not token

	// token

	jwt := user.Use(ginx.JwtAuth())
	jwt.GET("info", c.Info)      // 用户信息
	jwt.GET("list", c.List)      // 用户列表
	jwt.POST("update", c.Update) // 修改用户信息
}
