package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xgorm"
	"github.com/laixhe/gonet/xlog"

	"webapi/app/models"
	"webapi/protocol/gen/pbuser"
)

func (s *Service) UserInfo(c *gin.Context, req *pbuser.InfoRequest) (*pbuser.InfoResponse, xerror.IError) {
	uid := xgin.ContextUid64(c)
	if uid == 0 {
		return nil, xgin.IErrorAuthInvalid(nil)
	}
	//
	user, err := s.data.User.FirstID(uid)
	if err != nil {
		if xgorm.IsRecordNotFound(err) {
			return nil, xgin.IErrorAuthInvalid(err)
		}
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	return &pbuser.InfoResponse{
		User: &pbuser.User{
			Uid:       user.ID,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}, nil
}

func (s *Service) UserList(c *gin.Context, req *pbuser.ListRequest) (*pbuser.ListResponse, xerror.IError) {
	users, total, err := s.data.User.List(int(req.Size), int(req.Page))
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	resp := &pbuser.ListResponse{
		List:  make([]*pbuser.User, 0, len(users)),
		Total: int32(total),
		Size:  req.Size,
		Page:  req.Page,
	}
	for _, user := range users {
		resp.List = append(resp.List, &pbuser.User{
			Uid:       user.ID,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		})
	}
	//
	return resp, nil
}

func (s *Service) UserUpdate(c *gin.Context, req *pbuser.UpdateRequest) (*pbuser.UpdateResponse, xerror.IError) {
	uid := xgin.ContextUid64(c)
	if uid == 0 {
		return nil, xgin.IErrorAuthInvalid(nil)
	}
	//
	user, err := s.data.User.FirstUname(req.Uname)
	if err != nil && !xgorm.IsRecordNotFound(err) {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	if err == nil {
		if user.ID == uid {
			return &pbuser.UpdateResponse{}, nil
		}
		return nil, xgin.IErrorParse(errors.New("用户名已存在！"))
	}
	//
	user = models.User{
		ID:    uid,
		Uname: req.Uname,
	}
	err = s.data.User.Update(&user)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	return &pbuser.UpdateResponse{}, nil
}
