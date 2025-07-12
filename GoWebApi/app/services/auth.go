package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/xcrypto"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"webapi/app/entity"
	"webapi/app/model"
	"webapi/app/model/dao"
	"webapi/core"
	"webapi/core/middlewares"
)

// Auth 鉴权相关
type Auth struct {
	server *core.Server
	dao    *dao.Dao
}

func NewAuth(server *core.Server, modelDao *dao.Dao) *Auth {
	return &Auth{
		server: server,
		dao:    modelDao,
	}
}

// Register 注册
func (s *Auth) Register(ctx *fiber.Ctx, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error) {
	password, err := xcrypto.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, err
	}
	_, err = s.dao.GetUserByEmail(ctx.UserContext(), req.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "邮箱已存在")
	}
	user := &model.User{
		TypeId:    model.UserTypeOrdinary,
		Account:   xid.New().String(),
		Mobile:    "",
		Nickname:  req.Nickname,
		Email:     req.Email,
		Password:  password,
		AvatarUrl: "",
		Sex:       model.SexUnknown,
		States:    model.StateNormal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.dao.CreateUser(ctx.UserContext(), user)
	if err != nil {
		return nil, err
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthRegisterResponse{
		Token: token,
		User: &entity.User{
			Uid:       user.ID,
			TypeId:    user.TypeId,
			Account:   user.Account,
			Mobile:    user.Mobile,
			Email:     user.Email,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Sex:       user.Sex,
			States:    user.States,
		},
	}, nil
}

// Login 登录
func (s *Auth) Login(ctx *fiber.Ctx, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error) {
	user, err := s.dao.GetUserByEmail(ctx.UserContext(), req.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "邮箱或密码错误")
	}
	if !xcrypto.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "邮箱或密码错误")
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthLoginResponse{
		Token: token,
		User: &entity.User{
			Uid:       user.ID,
			TypeId:    user.TypeId,
			Account:   user.Account,
			Mobile:    user.Mobile,
			Email:     user.Email,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Sex:       user.Sex,
			States:    user.States,
		},
	}, nil
}

// 刷新Jwt
func (s *Auth) Refresh(ctx *fiber.Ctx, req *entity.AuthRefreshRequest) (*entity.AuthRefreshResponse, error) {
	user, err := s.dao.GetUserByID(ctx.UserContext(), req.Uid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthRefreshResponse{
		Token: token,
		User: &entity.User{
			Uid:       user.ID,
			TypeId:    user.TypeId,
			Account:   user.Account,
			Mobile:    user.Mobile,
			Email:     user.Email,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			Sex:       user.Sex,
			States:    user.States,
		},
	}, nil
}
