package services

import (
	"errors"
	"fmt"
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
	// 先查邮箱是否已注册, 避免无效的 bcrypt 计算;
	// 并发下的重复注册由 user.email 唯一索引兜底 (CreateUser 失败时按 ErrDuplicatedKey 处理)
	user := &models.User{}
	err := s.server.Gorm(ctx.Context()).Select("id").First(user, "email = ?", req.Email).Error
	if err == nil {
		return nil, xfiber.ParamError("邮箱已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("register: query email: %w", err)
	}

	password, err := crypto.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}
	user = &models.User{
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, xfiber.ParamError("邮箱已存在")
		}
		return nil, fmt.Errorf("register: create user: %w", err)
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, fmt.Errorf("register: generate token: %w", err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xfiber.ParamError("邮箱或密码错误")
		}
		return nil, fmt.Errorf("login: query user by email: %w", err)
	}
	// 封禁账号与密码错误返回同一提示, 避免暴露账号状态 (可被探测)
	if user.States != models.UserStateNormal {
		return nil, xfiber.ParamError("邮箱或密码错误")
	}
	if !crypto.BcryptPasswordCheck(req.Password, user.Password) {
		return nil, xfiber.ParamError("邮箱或密码错误")
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, fmt.Errorf("login: generate token: %w", err)
	}
	return &entity.AuthLoginResponse{
		Token: token,
		User:  entity.NewUserFromModel(user, "", ""),
	}, nil
}

// Refresh 刷新Jwt
func (s *Auth) Refresh(ctx fiber.Ctx, req *entity.AuthRefreshRequest) (*entity.AuthRefreshResponse, error) {
	// 只需要 states 与响应字段, 排除 password 减少不必要的列传输
	user := &models.User{}
	err := s.server.Gorm(ctx.Context()).
		Select(models.UserColumnsNoPassword).
		Where("id", req.Uid).
		Take(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xfiber.AuthorizedError()
		}
		return nil, fmt.Errorf("refresh: query user by id: %w", err)
	}
	if user.States != models.UserStateNormal {
		return nil, xfiber.AuthorizedError()
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, fmt.Errorf("refresh: generate token: %w", err)
	}
	return &entity.AuthRefreshResponse{
		Token: token,
		User:  entity.NewUserFromModel(user, "", ""),
	}, nil
}
