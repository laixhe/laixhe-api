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
	//g := gin.Default()
	g := gin.New()
	{
		// 中间件
		g.Use(xgin.SetRequestID()) // 设置请求ID
		g.Use(xgin.Cors())         // 跨域
		g.Use(xgin.Logger())       // 日志
		g.Use(xgin.Recovery())     // 出现 panic 恢复正常
		// 分组
		apiRouter := g.Group("api")
		{
			AuthRouter(apiRouter) // 鉴权相关
			UserRouter(apiRouter) // 用户相关
		}
		// doc 接口文档
		g.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	}
	return g
}
