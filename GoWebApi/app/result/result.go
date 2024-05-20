package result

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	"webapi/core/errorx"
	"webapi/core/utils"
	pbCode "webapi/profile/gen/code"
)

// Result 响应请求的公共模型
type Result struct {
	Code pbCode.ECode `json:"code"` // 响应码
	Msg  string       `json:"msg"`  // 响应信息
	Data any          `json:"data"` // 数据
}

func Response(c *gin.Context, code pbCode.ECode, data any, err string) {
	c.JSON(http.StatusOK, &Result{
		Code: code,
		Msg:  err,
		Data: data,
	})
}

func ResponseCode(c *gin.Context, code pbCode.ECode, err string) {
	c.JSON(http.StatusOK, &Result{
		Code: code,
		Msg:  err,
	})
}

func ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &Result{
		Code: pbCode.ECode_Success,
		Data: data,
	})
}

func ResponseError(c *gin.Context, err error) {
	if e, ok := err.(*errorx.Error); ok {
		c.JSON(http.StatusOK, &Result{
			Code: e.Code,
			Msg:  e.Error(),
		})
		return
	}
	if e, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusOK, &Result{
			Code: pbCode.ECode_Param,
			Msg:  utils.TranslatorErrorString(e),
		})
		return
	}
	if e, ok := err.(*json.UnmarshalTypeError); ok {
		c.JSON(http.StatusOK, &Result{
			Code: pbCode.ECode_Param,
			Msg:  e.Error(),
		})
		return
	}
	if e, ok := err.(*json.SyntaxError); ok {
		c.JSON(http.StatusOK, &Result{
			Code: pbCode.ECode_Param,
			Msg:  e.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Result{
		Code: pbCode.ECode_Unknown,
		Msg:  err.Error(),
	})
}
