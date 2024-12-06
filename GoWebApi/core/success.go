package core

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"
)

func JSONSuccess(c *gin.Context, data any) {
	xgin.SuccessJSON(c, data)
}
