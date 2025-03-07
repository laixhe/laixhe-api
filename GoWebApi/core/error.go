package core

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"

	"webapi/protocol/gen/ecode"
)

func NewError(code ecode.ECode, err error) xerror.IError {
	if err == nil {
		err = errors.New(code.String())
	}
	return xerror.New(int32(code), err)
}

func ErrorService(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Service.String())
	}
	return xerror.New(int32(ecode.ECode_Service), err)
}

func ErrorParse(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Parse.String())
	}
	return xerror.New(int32(ecode.ECode_Parse), err)
}

func ErrorEncode(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Encode.String())
	}
	return xerror.New(int32(ecode.ECode_Encode), err)
}

func ErrorParam(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Param.String())
	}
	return xerror.New(int32(ecode.ECode_Param), err)
}

func ErrorTip(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Tip.String())
	}
	return xerror.New(int32(ecode.ECode_Tip), err)
}

func ErrorRepeat(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_Repeat.String())
	}
	return xerror.New(int32(ecode.ECode_Repeat), err)
}

func ErrorAuthInvalid(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_AuthInvalid.String())
	}
	return xerror.New(int32(ecode.ECode_AuthInvalid), err)
}

func ErrorAuthExpire(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_AuthExpire.String())
	}
	return xerror.New(int32(ecode.ECode_AuthExpire), err)
}

func ErrorAuthUserError(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_AuthUserError.String())
	}
	return xerror.New(int32(ecode.ECode_AuthUserError), err)
}

func ErrorUserNotExist(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_UserNotExist.String())
	}
	return xerror.New(int32(ecode.ECode_UserNotExist), err)
}

func ErrorUserExist(err error) xerror.IError {
	if err == nil {
		err = errors.New(ecode.ECode_UserExist.String())
	}
	return xerror.New(int32(ecode.ECode_UserExist), err)
}

func NewErrorStr(code ecode.ECode, errStr string) xerror.IError {
	if errStr == "" {
		errStr = code.String()
	}
	return xerror.NewStr(int32(code), errStr)
}

func ErrorParamStr(errStr string) xerror.IError {
	if errStr == "" {
		errStr = ecode.ECode_Param.String()
	}
	return xerror.NewStr(int32(ecode.ECode_Param), errStr)
}

func JSONError(c *gin.Context, err xerror.IError) {
	xgin.ErrorJSON(c, err)
}

func JSONErrorService(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorService(err))
}

func JSONErrorParse(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorParse(err))
}

func JSONErrorEncode(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorEncode(err))
}

func JSONErrorParam(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorParam(err))
}

func JSONErrorTip(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorTip(err))
}

func JSONErrorRepeat(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorRepeat(err))
}

func JSONErrorAuthInvalid(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorAuthInvalid(err))
}

func JSONErrorAuthExpire(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorAuthExpire(err))
}

func JSONErrorAuthUserError(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorAuthUserError(err))
}

func JSONErrorUserNotExist(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorUserNotExist(err))
}

func JSONErrorUserExist(c *gin.Context, err error) {
	xgin.ErrorJSON(c, ErrorUserExist(err))
}

func JSONErrorParamStr(c *gin.Context, errStr string) {
	xgin.ErrorJSON(c, ErrorParamStr(errStr))
}
