package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/errorx"
	"github.com/laixhe/gonet/ginx"
	"github.com/laixhe/gonet/jwtx"
	"github.com/laixhe/gonet/logx"
	"github.com/laixhe/gonet/utils"
	"github.com/rs/xid"
	"gorm.io/gorm"

	pbAuth "webapi/api/gen/auth"
	"webapi/api/gen/code"
	pbUser "webapi/api/gen/user"
	"webapi/app/models"
	"webapi/core"
)

// AuthRegister 注册
func (s *Service) AuthRegister(c *gin.Context, req *pbAuth.RegisterRequest) (*pbAuth.RegisterResponse, error) {
	u, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("FirstEmail %v", err)
			return nil, errorx.ServiceError(err)
		}
	}
	if u.Uid > 0 {
		return nil, errorx.ParamError(nil)
	}

	password, err := utils.BcryptPasswordHash(req.Password)
	if err != nil {
		logx.Errorf("FirstEmail %v", err)
		return nil, errorx.ServiceError(err)
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
		logx.Errorf("Create %v", err)
		return nil, errorx.ServiceError(err)
	}

	token, err := jwtx.GenToken(core.Config().Jwt, user.Uid, xid.New().String())
	if err != nil {
		logx.Errorf("GenToken %v", err)
		return nil, errorx.ServiceError(err)
	}
	resp := &pbAuth.RegisterResponse{
		Token: token,
		User: &pbUser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

// AuthLogin 登录
func (s *Service) AuthLogin(c *gin.Context, req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error) {
	user, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New(int32(code.Code_AuthUserError), nil)
		}
		logx.Errorf("FirstEmail %v", err)
		return nil, errorx.ServiceError(err)
	}
	if !utils.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, errorx.New(int32(code.Code_AuthUserError), nil)
	}
	token, err := jwtx.GenToken(core.Config().Jwt, user.Uid, xid.New().String())
	if err != nil {
		logx.Errorf("GenToken %v", err)
		return nil, errorx.ServiceError(err)
	}
	resp := &pbAuth.LoginResponse{
		Token: token,
		User: &pbUser.User{
			Uid:       user.Uid,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}
	return resp, nil
}

// AuthRefresh 刷新Jwt
func (s *Service) AuthRefresh(c *gin.Context, req *pbAuth.RefreshRequest) (*pbAuth.RefreshResponse, error) {
	uid, err := ginx.ContextUid(c)
	if err != nil {
		return nil, err
	}
	token, err := jwtx.GenToken(core.Config().Jwt, uid, xid.New().String())
	if err != nil {
		logx.Errorf("GenToken %v", err)
		return nil, errorx.ServiceError(err)
	}
	resp := &pbAuth.RefreshResponse{
		Token: token,
	}
	return resp, nil
}
