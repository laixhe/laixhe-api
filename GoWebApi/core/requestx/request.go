package requestx

import (
	"github.com/gin-gonic/gin"

	"webapi/api/gen/enum/eapp"
)

// 头部 平台

const PlatformHeaderKey = "platform"

func GetPlatform(c *gin.Context) eapp.Platform {
	platform := c.Request.Header.Get(PlatformHeaderKey)
	if platform != "" {
		value := eapp.Platform_value[platform]
		if value != 0 {
			return eapp.Platform(value)
		}
	}
	return eapp.Platform_unknown
}
