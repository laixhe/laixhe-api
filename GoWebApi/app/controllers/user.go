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

	result.ResponseSuccess(c, resp)
}

// List 用户列表
//
// @Summary	用户列表
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  false "Bearer token令牌"
// @Param    size          query     string  false "每页页数(数量)"
// @Param    page          query     string  false "当前页数"
// @Success  200           {object}  user.ListResponse
// @Success  400           {object}  result.Result
// @Router   /api/user/list [get]
func (u *User) List(c *gin.Context) {
	req := &pbUser.ListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		logx.Infof("req:%s", req)
		result.ResponseError(c, err)
		return
	}
	logx.Infof("req:%s", req)
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	resp, err := u.service.UserList(c, req)
	if err != nil {
		result.ResponseError(c, err)
		return
	}

	result.ResponseSuccess(c, resp)
}
