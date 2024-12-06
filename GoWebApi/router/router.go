package router

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/proto/gen/config/clog"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xgin/xvalidator"
	"github.com/laixhe/gonet/xlog"
	"github.com/rs/xid"
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
	if err := xvalidator.ValidatorTranslator(xvalidator.Zh); err != nil {
		xlog.Errorf("gin set translator error:%v", err)
		return nil
	}

	//r := gin.Default()
	r := gin.New()
	// 中间件
	r.Use(requestid.New(requestid.WithGenerator(func() string {
		return xid.New().String()
	})))                   // 请求ID
	r.Use(xgin.Cors())     // 跨域
	r.Use(xgin.Logger())   // 日志
	r.Use(xgin.Recovery()) // 出现 panic 恢复正常

	api := r.Group("/api")
	AuthRouter(api) // 鉴权相关
	UserRouter(api) // 用户相关

	// doc 接口文档
	r.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	return r
}
