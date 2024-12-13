package controllers

import (
	"time"

	"github.com/gin-gonic/gin"

	"webapi/app/services"
	"webapi/core"
	"webapi/protocol/gen/pbuser"
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
// @Success  200    {object}  pbuser.InfoRequest
// @Router   /api/user/info [get]
func (u *User) Info(c *gin.Context) {
	req := &pbuser.InfoRequest{}
	resp, err := u.service.UserInfo(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	core.JSONSuccess(c, resp)
}

// List 用户列表
//
// @Summary	用户列表
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  false "Bearer token令牌"
// @Param    size          query     string  false "每页页数(数量)"
// @Param    page          query     string  false "当前页数"
// @Success  200           {object}  pbuser.ListResponse
// @Router   /api/user/list [get]
func (u *User) List(c *gin.Context) {
	req := &pbuser.ListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		core.JSONErrorParse(c, err)
		return
	}
	//
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	//
	resp, err := u.service.UserList(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	core.JSONSuccess(c, resp)
}

// Update 修改用户信息
//
// @Summary	修改用户信息
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  false "Bearer token令牌"
// @Param    body          body      pbuser.UpdateRequest   ture "请求body参数"
// @Success  200           {object}  pbuser.UpdateResponse
// @Router   /api/user/update [post]
func (u *User) Update(c *gin.Context) {
	req := &pbuser.UpdateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		core.JSONErrorParse(c, err)
		return
	}
	//
	if _, err := time.ParseInLocation(time.DateTime, req.LoginAt, time.Local); err != nil {
		core.JSONErrorParamStr(c, "登录时间格式不对！")
		return
	}
	//
	resp, err := u.service.UserUpdate(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	core.JSONSuccess(c, resp)
}
