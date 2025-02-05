package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xcrypto"
	"github.com/laixhe/gonet/xerror"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xgorm"
	"github.com/laixhe/gonet/xjwt"
	"github.com/laixhe/gonet/xlog"
	"github.com/rs/xid"

	"webapi/app/models"
	"webapi/core"
	"webapi/protocol/gen/pbauth"
	"webapi/protocol/gen/pbuser"
)

// AuthRegister 注册
func (s *Service) AuthRegister(c *gin.Context, req *pbauth.RegisterRequest) (*pbauth.RegisterResponse, xerror.IError) {
	u, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if !xgorm.IsRecordNotFound(err) {
			xlog.Error(err.Error(), xgin.ZapField(c)...)
			return nil, xgin.IErrorServer(err)
		}
	}
	if u.ID > 0 {
		return nil, xgin.IErrorParse(nil)
	}
	//
	password, err := xcrypto.BcryptPasswordHash(req.Password)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	user := &models.User{
		Password: password,
		Email:    req.Email,
		Uname:    req.Uname,
		Age:      req.Age,
		Score:    0,
	}
	err = s.data.User.Create(user)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	claims := &xjwt.CustomClaims{Uid: int(user.ID)}
	claims.ID = xid.New().String()
	token, err := xjwt.GenToken(core.Config().Jwt, claims)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	return &pbauth.RegisterResponse{
		Token: token,
		User: &pbuser.User{
			Uid:       user.ID,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}, nil
}

// AuthLogin 登录
func (s *Service) AuthLogin(c *gin.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, xerror.IError) {
	user, err := s.data.User.FirstEmail(req.Email)
	if err != nil {
		if xgorm.IsRecordNotFound(err) {
			return nil, xgin.IErrorAuthInvalid(nil)
		}
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	if !xcrypto.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, xgin.IErrorAuthInvalid(nil)
	}
	claims := &xjwt.CustomClaims{Uid: int(user.ID)}
	claims.ID = xid.New().String()
	token, err := xjwt.GenToken(core.Config().Jwt, claims)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	return &pbauth.LoginResponse{
		Token: token,
		User: &pbuser.User{
			Uid:       user.ID,
			Uname:     user.Uname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
		},
	}, nil
}

// AuthRefresh 刷新Jwt
func (s *Service) AuthRefresh(c *gin.Context, req *pbauth.RefreshRequest) (*pbauth.RefreshResponse, xerror.IError) {
	uid := xgin.ContextUid(c)
	if uid == 0 {
		return nil, xgin.IErrorAuthInvalid(nil)
	}
	//
	claims := &xjwt.CustomClaims{Uid: uid}
	claims.ID = xid.New().String()
	token, err := xjwt.GenToken(core.Config().Jwt, claims)
	if err != nil {
		xlog.Error(err.Error(), xgin.ZapField(c)...)
		return nil, xgin.IErrorServer(err)
	}
	//
	return &pbauth.RefreshResponse{
		Token: token,
	}, nil
}
