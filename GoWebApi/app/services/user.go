package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webapi/app/middleware"
	"webapi/core/errorx"
	"webapi/core/logx"
	pbCode "webapi/profile/gen/code"
	pbUser "webapi/profile/gen/user"
)

func (s *Service) UserInfo(c *gin.Context, req *pbUser.InfoRequest) (*pbUser.InfoResponse, error) {
	uid, errx := middleware.Uid(c)
	if errx != nil {
		return nil, errx
	}

	user, err := s.data.User.FirstID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewError(pbCode.ECode_AuthNotLogin, err)
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
	users, err := s.data.User.List()
	if err != nil {
		logx.Errorf("UserList %v", err)
		return nil, errorx.ServiceError(err)
	}

	resp := &pbUser.ListResponse{
		Users: make([]*pbUser.User, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &pbUser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		})
	}
	return resp, nil
}
