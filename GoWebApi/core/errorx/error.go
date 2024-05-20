package errorx

import pbCode "webapi/profile/gen/code"

type Error struct {
	Code pbCode.ECode
	Err  error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Code.String()
	}
	return e.Err.Error()
}

func NewError(code pbCode.ECode, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func UnknownError(err error) *Error {
	return &Error{
		Code: pbCode.ECode_Unknown,
		Err:  err,
	}
}

func ServiceError(err error) *Error {
	return &Error{
		Code: pbCode.ECode_Service,
		Err:  err,
	}
}

func ParamError(err error) *Error {
	return &Error{
		Code: pbCode.ECode_Param,
		Err:  err,
	}
}
