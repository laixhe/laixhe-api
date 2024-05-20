package controllers

import (
	"github.com/gin-gonic/gin"

	"webapi/app/result"
	"webapi/app/services"
	"webapi/core/logx"
	pbUser "webapi/profile/gen/user"
)

type User struct {
	service *services.Service
}

func NewUser() *User {
	return &User{
		service: services.NewService(),
	}
}

// Info 用户信息
//
// @Summary	用户信息
// @Accept   json
// @Produce  json
// @Param Authorization header string false "Bearer token令牌"
// @Success  200    {object}  user.InfoResponse
// @Success  400    {object}  result.Result
// @Router   /api/user/info [get]
func (u *User) Info(c *gin.Context) {
	req := &pbUser.InfoRequest{}
	resp, err := u.service.UserInfo(c, req)
	if err != nil {
		result.ResponseError(c, err)
		return
	}

	logx.Infof("resp:%s", resp)
	result.ResponseSuccess(c, resp)
}

// List 用户列表
//
// @Summary	用户列表
// @Accept   json
// @Produce  json
// @Param Authorization header string false "Bearer token令牌"
// @Success  200    {object}  user.ListResponse
// @Success  400    {object}  result.Result
// @Router   /api/user/list [get]
func (u *User) List(c *gin.Context) {
	req := &pbUser.ListRequest{}
	resp, err := u.service.UserList(c, req)
	if err != nil {
		result.ResponseError(c, err)
		return
	}

	logx.Infof("resp:%s", resp)
	result.ResponseSuccess(c, resp)
}
