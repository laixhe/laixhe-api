package responsex

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"

	"webapi/api/gen/ecode"
	"webapi/core/errorx"
	"webapi/core/utils"
)

// ResponseModel 响应请求的公共模型
type ResponseModel struct {
	Code ecode.ECode `json:"code"`           // 响应码
	Msg  string      `json:"msg"`            // 响应信息
	Data any         `json:"data,omitempty"` // 数据
}

func ResponseError(err error) *ResponseModel {
	var ex *errorx.Error
	if errors.As(err, &ex) {
		return &ResponseModel{
			Code: ex.Code,
			Msg:  ex.Error(),
		}
	}
	var ev validator.ValidationErrors
	if errors.As(err, &ev) {
		return &ResponseModel{
			Code: ecode.ECode_Param,
			Msg:  utils.TranslatorErrorString(ev),
		}
	}
	var ejut *json.UnmarshalTypeError
	if errors.As(err, &ejut) {
		return &ResponseModel{
			Code: ecode.ECode_Param,
			Msg:  ejut.Error(),
		}
	}
	var ejse *json.SyntaxError
	if errors.As(err, &ejse) {
		return &ResponseModel{
			Code: ecode.ECode_Param,
			Msg:  ejse.Error(),
		}
	}
	return &ResponseModel{
		Code: ecode.ECode_Service,
		Msg:  err.Error(),
	}
}
