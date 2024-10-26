package router

import (
	"github.com/gin-gonic/gin"

	"webapi/app/controllers"
	"webapi/core/ginx"
)

// AuthRouter 鉴权相关
func AuthRouter(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	c := controllers.NewAuth()

	// not token

	auth.POST("/register", c.Register) // 注册
	auth.POST("/login", c.Login)       // 登陆

	// token

	jwt := auth.Use(ginx.JwtAuth())
	jwt.POST("/refresh", c.Refresh) // 刷新Jwt
}
