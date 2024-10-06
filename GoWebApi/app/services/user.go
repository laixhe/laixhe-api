package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"

	pbUser "webapi/api/gen/user"
	"webapi/app/models"
	"webapi/core/errorx"
	"webapi/core/ginx"
	"webapi/core/logx"
)

func (s *Service) UserInfo(c *gin.Context, req *pbUser.InfoRequest) (*pbUser.InfoResponse, error) {
	uid, err := ginx.ContextUid(c)
	if err != nil {
		return nil, err
	}

	user, err := s.data.User.FirstID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.AuthInvalidError(err)
		}
		logx.Errorf("UserInfo %v", err)
		return nil, errorx.ServiceError(err)
	}
	resp := &pbUser.InfoResponse{
		User: &pbUser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

func (s *Service) UserList(c *gin.Context, req *pbUser.ListRequest) (*pbUser.ListResponse, error) {
	users, total, err := s.data.User.List(int(req.Size), int(req.Page))
	if err != nil {
		logx.Errorf("UserList %v", err)
		return nil, errorx.ServiceError(err)
	}

	resp := &pbUser.ListResponse{
		List:  make([]*pbUser.User, 0, len(users)),
		Total: int32(total),
		Size:  req.Size,
		Page:  req.Page,
	}
	for _, user := range users {
		resp.List = append(resp.List, &pbUser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		})
	}
	return resp, nil
}

func (s *Service) UserUpdate(c *gin.Context, req *pbUser.UpdateRequest) (*pbUser.UpdateResponse, error) {
	uid, err := ginx.ContextUid(c)
	if err != nil {
		return nil, err
	}

	user, err := s.data.User.FirstUname(req.Uname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.Errorf("UserUpdate %v", err)
		return nil, errorx.ServiceError(err)
	}
	if err == nil {
		if user.Uid == uid {
			return &pbUser.UpdateResponse{}, nil
		}
		return nil, errorx.ParamError(errors.New("用户名已存在！"))
	}

	user = models.User{
		Uid:     uid,
		Uname:   req.Uname,
		LoginAt: carbon.Parse(req.LoginAt).StdTime(),
	}
	err = s.data.User.Update(&user)
	if err != nil {
		logx.Errorf("Update %v", err)
		return nil, errorx.ServiceError(err)
	}

	return &pbUser.UpdateResponse{}, nil
}
