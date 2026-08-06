package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/crypto"
	"github.com/laixhe/gonet/jwt"
	"github.com/laixhe/gonet/xfiber"
	"github.com/rs/xid"
	"gorm.io/gorm"

	"webapi/app/entity"
	"webapi/app/models"
	"webapi/core"
	"webapi/core/middlewares"
)

// Auth 鉴权相关
type Auth struct {
	server *core.Server
}

// NewAuth 创建 Auth 业务实例
func NewAuth(server *core.Server) *Auth {
	return &Auth{
		server: server,
	}
}

// Register 注册
func (s *Auth) Register(ctx fiber.Ctx, req *entity.AuthRegisterRequest) (*entity.AuthRegisterResponse, error) {
	// 先检查邮箱是否已注册，避免无效的 bcrypt 计算
	{
		user := &models.User{}
		err := s.server.Gorm(ctx.Context()).Select("id").First(user, "email = ?", req.Email).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			return nil, xfiber.ParamError("邮箱已存在")
		}
	}
	password, err := crypto.BcryptPasswordHash(req.Password)
	user := &models.User{
		TypeId:    models.UserTypeOrdinary,
		Account:   xid.New().String(),
		Mobile:    "",
		Nickname:  req.Nickname,
		Email:     req.Email,
		Password:  password,
		AvatarUrl: "",
		Sex:       models.UserSexUnknown,
		States:    models.UserStateNormal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err = models.CreateUser(s.server.Gorm(ctx.Context()), user); err != nil {
		return nil, err
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthRegisterResponse{
		Token: token,
		User:  entity.NewUserFromModel(user, "", ""),
	}, nil
}

// Login 登录
func (s *Auth) Login(ctx fiber.Ctx, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error) {
	user := &models.User{}
	if err := s.server.Orm().FirstByField(ctx.Context(), user, "email", req.Email); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.ParamError("邮箱或密码错误")
	}
	if user.States != models.UserStateNormal {
		return nil, xfiber.AuthorizedError()
	}
	if !crypto.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, xfiber.ParamError("邮箱或密码错误")
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthLoginResponse{
		Token: token,
		User:  entity.NewUserFromModel(user, "", ""),
	}, nil
}

// Refresh 刷新Jwt
func (s *Auth) Refresh(ctx fiber.Ctx, req *entity.AuthRefreshRequest) (*entity.AuthRefreshResponse, error) {
	user := &models.User{}
	if err := s.server.Orm().GetById(ctx.Context(), user, req.Uid); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.AuthorizedError()
	}
	if user.States != models.UserStateNormal {
		return nil, xfiber.AuthorizedError()
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, err
	}
	return &entity.AuthRefreshResponse{
		Token: token,
		User:  entity.NewUserFromModel(user, "", ""),
	}, nil
}
