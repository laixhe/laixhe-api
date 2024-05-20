package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webapi/app/middleware"
	"webapi/app/models"
	"webapi/core/config"
	"webapi/core/errorx"
	"webapi/core/jwtx"
	"webapi/core/logx"
	"webapi/core/utils"
	pbAuth "webapi/profile/gen/auth"
	pbCode "webapi/profile/gen/code"
	pbUser "webapi/profile/gen/user"
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
		return nil, errorx.NewError(pbCode.ECode_EmailExist, nil)
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

	token, err := jwtx.GenToken(user.Uid, config.Get().Jwt)
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
			return nil, errorx.NewError(pbCode.ECode_AuthUserError, nil)
		}
		logx.Errorf("FirstEmail %v", err)
		return nil, errorx.ServiceError(err)
	}
	if !utils.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, errorx.NewError(pbCode.ECode_AuthUserError, err)
	}
	token, err := jwtx.GenToken(user.Uid, config.Get().Jwt)
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
	uid, errx := middleware.Uid(c)
	if errx != nil {
		return nil, errx
	}
	token, err := jwtx.GenToken(uid, config.Get().Jwt)
	if err != nil {
		logx.Errorf("GenToken %v", err)
		return nil, errorx.ServiceError(err)
	}
	resp := &pbAuth.RefreshResponse{
		Token: token,
	}
	return resp, nil
}
