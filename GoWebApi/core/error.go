package core

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"

	"webapi/api/gen/ecode"
)

func NewError(code ecode.ECode, err error) xerror.IError {
	return xerror.New(int32(code), err)
}

func NewErrorStr(code ecode.ECode, errStr string) xerror.IError {
	return xerror.NewStr(int32(code), errStr)
}

func ErrorService(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Service), err)
}

func ErrorParse(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Parse), err)
}

func ErrorEncode(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Encode), err)
}

func ErrorParam(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Param), err)
}

func ErrorParamStr(errStr string) xerror.IError {
	return xerror.NewStr(int32(ecode.ECode_Param), errStr)
}

func ErrorTip(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Tip), err)
}

func ErrorRepeat(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_Repeat), err)
}

func ErrorAuthInvalid(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_AuthInvalid), err)
}

func ErrorAuthExpire(err error) xerror.IError {
	return xerror.New(int32(ecode.ECode_AuthExpire), err)
}

func JSONError(c *gin.Context, err xerror.IError) {
	xgin.ErrorJSON(c, err)
}

func JSONErrorService(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Service), err))
}

func JSONErrorParse(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Parse), err))
}

func JSONErrorEncode(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Encode), err))
}

func JSONErrorParam(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Param), err))
}

func JSONErrorParamStr(c *gin.Context, errStr string) {
	xgin.ErrorJSON(c, xerror.NewStr(int32(ecode.ECode_Param), errStr))
}

func JSONErrorTip(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Tip), err))
}

func JSONErrorRepeat(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_Repeat), err))
}

func JSONErrorAuthInvalid(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_AuthInvalid), err))
}

func JSONErrorAuthExpire(c *gin.Context, err error) {
	xgin.ErrorJSON(c, xerror.New(int32(ecode.ECode_AuthExpire), err))
}
