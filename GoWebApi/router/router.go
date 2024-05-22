package router

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"

	"webapi/core/config"
	"webapi/core/logx"
	"webapi/core/utils"
)

// Router gin 路由
func Router() *gin.Engine {
	if config.Get().Log.Level == logx.LogLevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Validator(表单验证)多语言提示文本
	if err := utils.ValidatorTranslator("zh"); err != nil {
		logx.Errorf("gin set translator error:%v", err)
		return nil
	}

	//r := gin.Default()
	r := gin.New()
	// 中间件
	r.Use(requestid.New(requestid.WithGenerator(func() string {
		return xid.New().String()
	})))                 // 请求ID
	r.Use(Cors())        // 跨域
	r.Use(GinLogger())   // 日志
	r.Use(GinRecovery()) // 出现 panic 恢复正常

	api := r.Group("/api")
	AuthRouter(api) // 鉴权相关-登录、注册
	UserRouter(api) // 用户相关

	// doc 接口文档
	r.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	return r
}
