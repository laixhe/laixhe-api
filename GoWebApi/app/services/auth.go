package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
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
	// email 为唯一索引 (与 webapi.sql 一致), 先查后插 + 数据库唯一约束双重防重。
	// 只取 id 列避免拉取整行 (查询风格说明见 README「核心封装速查-查询风格」)
	user := &models.User{}
	err := s.server.Gorm(ctx.Context()).Select("id").First(user, "email = ?", req.Email).Error
	if err == nil {
		return nil, xfiber.ParamError("邮箱已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("register: query email: %w", err)
	}

	// bcrypt cost=10 为 CPU 密集计算 (单次约 50-100ms), 提交到 BcryptPool worker 池执行,
	// 避免占用请求 goroutine 所在线程 (对齐 Rust 版 spawn_blocking 思路, 见 core/bcrypt_pool.go)。
	// 池大小 = GOMAXPROCS, 并发超过池大小时请求排队等待 (背压),
	// 配合接口限流 (默认 1000 次/分钟/IP) + 先查后插, 将影响控制在教学规模内。
	password, err := s.server.Bcrypt().Hash(req.Password)
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
		// 唯一键: account(xid)、email、各关联表 uid, 冲突仅在并发注册同邮箱等极端情况出现
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, xfiber.ParamError("注册失败，请稍后再试")
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
		User:  *entity.NewUserFromModel(user, "", ""),
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
	// bcrypt verify 与 hash 同为 CPU 密集计算, 同样提交到 BcryptPool worker 池执行 (见 Register 处注释)
	if !s.server.Bcrypt().Check(req.Password, user.Password) {
		return nil, xfiber.ParamError("邮箱或密码错误")
	}
	claims := middlewares.NewJwtClaims(user.ID, s.server.Config().Jwt.ExpireTime)
	token, err := jwt.GenToken(s.server.Config().Jwt, claims)
	if err != nil {
		return nil, fmt.Errorf("login: generate token: %w", err)
	}
	return &entity.AuthLoginResponse{
		Token: token,
		User:  *entity.NewUserFromModel(user, "", ""),
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
		User:  *entity.NewUserFromModel(user, "", ""),
	}, nil
}
