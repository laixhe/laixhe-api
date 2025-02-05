package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/protocol/gen/config/clog"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xlog"
	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"

	"webapi/core"
)

// Router gin 路由
func Router() *gin.Engine {
	if core.Config().Log.Level == clog.LevelType_debug.String() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// Validator(表单验证)多语言提示文本
	if err := xgin.ValidatorTranslator(xgin.Zh); err != nil {
		xlog.Errorf("gin set translator error:%v", err)
		return nil
	}
	//r := gin.Default()
	r := gin.New()
	{
		// 中间件
		r.Use(xgin.SetRequestID()) // 设置请求ID
		r.Use(xgin.Cors())         // 跨域
		r.Use(xgin.Logger())       // 日志
		r.Use(xgin.Recovery())     // 出现 panic 恢复正常
		// 分组
		apiRouter := r.Group("api")
		// 分组 v1
		v1Router := apiRouter.Group("v1")
		{
			AuthRouter(v1Router) // 鉴权相关
			UserRouter(v1Router) // 用户相关
		}
		// doc 接口文档
		r.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	}
	return r
}
