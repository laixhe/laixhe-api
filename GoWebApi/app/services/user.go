package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	carbon "github.com/dromara/carbon/v2"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xlog"
	"gorm.io/gorm"

	"webapi/api/gen/pbuser"
	"webapi/app/models"
)

func (s *Service) UserInfo(c *gin.Context, req *pbuser.InfoRequest) (*pbuser.InfoResponse, error) {
	uid, err := xgin.ContextUid(c)
	if err != nil {
		return nil, err
	}

	user, err := s.data.User.FirstID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.AuthInvalidError(err)
		}
		xlog.Errorf("UserInfo %v", err)
		return nil, xerror.ServiceError(err)
	}
	resp := &pbuser.InfoResponse{
		User: &pbuser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

func (s *Service) UserList(c *gin.Context, req *pbuser.ListRequest) (*pbuser.ListResponse, error) {
	users, total, err := s.data.User.List(int(req.Size), int(req.Page))
	if err != nil {
		xlog.Errorf("UserList %v", err)
		return nil, xerror.ServiceError(err)
	}

	resp := &pbuser.ListResponse{
		List:  make([]*pbuser.User, 0, len(users)),
		Total: int32(total),
		Size:  req.Size,
		Page:  req.Page,
	}
	for _, user := range users {
		resp.List = append(resp.List, &pbuser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		})
	}
	return resp, nil
}

func (s *Service) UserUpdate(c *gin.Context, req *pbuser.UpdateRequest) (*pbuser.UpdateResponse, error) {
	uid, err := xgin.ContextUid(c)
	if err != nil {
		return nil, err
	}

	user, err := s.data.User.FirstUname(req.Uname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		xlog.Errorf("UserUpdate %v", err)
		return nil, xerror.ServiceError(err)
	}
	if err == nil {
		if user.Uid == uid {
			return &pbuser.UpdateResponse{}, nil
		}
		return nil, xerror.ParamErrorStr("用户名已存在！")
	}

	user = models.User{
		Uid:     uid,
		Uname:   req.Uname,
		LoginAt: carbon.Parse(req.LoginAt).StdTime(),
	}
	err = s.data.User.Update(&user)
	if err != nil {
		xlog.Errorf("Update %v", err)
		return nil, xerror.ServiceError(err)
	}

	return &pbuser.UpdateResponse{}, nil
}
