package core

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"
)

// GinZapField gin日志字段
func GinZapField(c *gin.Context) []zap.Field {
	return []zap.Field{
		zap.String(xgin.HeaderRequestID, xgin.GetRequestID(c)),
		zap.String("path", c.Request.URL.Path),
	}
}
