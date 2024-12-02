package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xjwt"
	"github.com/laixhe/gonet/xlog"
	"github.com/laixhe/gonet/xutil"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"webapi/api/gen/code"
	"webapi/api/gen/pbauth"
	"webapi/api/gen/pbuser"
	"webapi/app/models"
	"webapi/core"
)

// AuthRegister 注册
func (s *Service) AuthRegister(c *gin.Context, req *pbauth.RegisterRequest) (*pbauth.RegisterResponse, error) {
	u, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			xlog.Errorf("FirstEmail %v", err)
			return nil, xerror.ServiceError(err)
		}
	}
	if u.Uid > 0 {
		return nil, xerror.ParamError(nil)
	}

	password, err := xutil.BcryptPasswordHash(req.Password)
	if err != nil {
		xlog.Errorf("FirstEmail %v", err)
		return nil, xerror.ServiceError(err)
	}
	user := &models.User{
		Password: password,
		Email:    req.Email,
		Uname:    req.Uname,
		Age:      req.Age,
		Score:    0,
		LoginAt:  time.Now(),
	}
	err = s.data.User.Create(user)
	if err != nil {
		xlog.Errorf("Create %v", err)
		return nil, xerror.ServiceError(err)
	}

	token, err := xjwt.GenToken(core.Config().Jwt, user.Uid, xid.New().String())
	if err != nil {
		xlog.Errorf("GenToken %v", err)
		return nil, xerror.ServiceError(err)
	}
	resp := &pbauth.RegisterResponse{
		Token: token,
		User: &pbuser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

// AuthLogin 登录
func (s *Service) AuthLogin(c *gin.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	user, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.New(int32(code.Code_AuthUserError), nil)
		}
		xlog.Errorf("FirstEmail %v", err)
		return nil, xerror.ServiceError(err)
	}
	if !xutil.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, xerror.New(int32(code.Code_AuthUserError), nil)
	}
	token, err := xjwt.GenToken(core.Config().Jwt, user.Uid, xid.New().String())
	if err != nil {
		xlog.Errorf("GenToken %v", err)
		return nil, xerror.ServiceError(err)
	}
	resp := &pbauth.LoginResponse{
		Token: token,
		User: &pbuser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

// AuthRefresh 刷新Jwt
func (s *Service) AuthRefresh(c *gin.Context, req *pbauth.RefreshRequest) (*pbauth.RefreshResponse, error) {
	uid, err := xgin.ContextUid(c)
	if err != nil {
		return nil, err
	}
	token, err := xjwt.GenToken(core.Config().Jwt, uid, xid.New().String())
	if err != nil {
		xlog.Errorf("GenToken %v", err)
		return nil, xerror.ServiceError(err)
	}
	resp := &pbauth.RefreshResponse{
		Token: token,
	}
	return resp, nil
}
