package controllers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	pbUser "webapi/api/gen/user"
	"webapi/app/services"
	"webapi/core/errorx"
	"webapi/core/ginx"
	"webapi/core/logx"
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
// @Success  200    {object}  user.InfoRequest
// @Success  400    {object}  responsex.ResponseModel
// @Router   /api/user/info [get]
func (u *User) Info(c *gin.Context) {
	req := &pbUser.InfoRequest{}
	resp, err := u.service.UserInfo(c, req)
	if err != nil {
		ginx.Error(c, err)
		return
	}

	ginx.Success(c, resp)
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
// @Success  400           {object}  responsex.ResponseModel
// @Router   /api/user/list [get]
func (u *User) List(c *gin.Context) {
	req := &pbUser.ListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		logx.Infof("req:%s", req)
		ginx.Error(c, err)
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
		ginx.Error(c, err)
		return
	}

	ginx.Success(c, resp)
}

// Update 修改用户信息
//
// @Summary	修改用户信息
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  false "Bearer token令牌"
// @Param    body          body      user.UpdateRequest   ture "请求body参数"
// @Success  200           {object}  user.UpdateResponse
// @Success  400           {object}  responsex.ResponseModel
// @Router   /api/user/update [post]
func (u *User) Update(c *gin.Context) {
	req := &pbUser.UpdateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		logx.Infof("req:%s", req)
		ginx.Error(c, err)
		return
	}
	logx.Infof("req:%s", req)

	_, err := time.ParseInLocation(time.DateTime, req.LoginAt, time.Local)
	if err != nil {
		ginx.Error(c, errorx.ParamError(errors.New("登录时间格式不对！")))
		return
	}

	resp, err := u.service.UserUpdate(c, req)
	if err != nil {
		ginx.Error(c, err)
		return
	}

	ginx.Success(c, resp)
}
