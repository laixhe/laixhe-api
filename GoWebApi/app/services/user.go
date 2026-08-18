package services

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/db/gorm/orm"
	"github.com/laixhe/gonet/xfiber"
	"gorm.io/gorm"

	"webapi/app/entity"
	"webapi/app/models"
	"webapi/core"
)

// User 用户相关
type User struct {
	server *core.Server
}

// NewUser 创建 User 业务实例
func NewUser(server *core.Server) *User {
	return &User{
		server: server,
	}
}

// Update 更新用户信息
//
// 先查后改: 查询 (排除 password) 用于 states 校验与响应组装;
// UpdateUser 按非零字段更新, 返回的是更新后的预期值而非 DB 回读值。
func (s *User) Update(ctx fiber.Ctx, req *entity.UserUpdateRequest) (*entity.User, error) {
	user := &models.User{}
	err := s.server.Gorm(ctx.Context()).
		Select(models.UserColumnsNoPassword).
		Where("id", req.Uid).
		Take(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xfiber.ParamError("用户不存在")
		}
		return nil, fmt.Errorf("user update: query user by id: %w", err)
	}
	if user.States != models.UserStateNormal {
		return nil, xfiber.AuthorizedError()
	}
	resp := entity.NewUserFromModel(user, req.Nickname, req.AvatarUrl)
	uid := user.ID
	user = &models.User{
		ID:        uid,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
	}
	if err := models.UpdateUser(s.server.Gorm(ctx.Context()), user); err != nil {
		return nil, fmt.Errorf("user update: update user: %w", err)
	}
	return resp, nil
}

// Info 获取用户信息
func (s *User) Info(ctx fiber.Ctx, req *entity.UserInfoRequest) (*entity.User, error) {
	user := &models.User{}
	err := s.server.Gorm(ctx.Context()).
		Select(models.UserColumnsNoPassword).
		Where("id", req.Uid).
		Take(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xfiber.ParamError("用户不存在")
		}
		return nil, fmt.Errorf("user info: query user by id: %w", err)
	}
	return entity.NewUserFromModel(user, "", ""), nil
}

// List 获取用户列表
func (s *User) List(ctx fiber.Ctx, req *entity.UserListRequest) (*entity.UserListResponse, error) {
	limit, offset := orm.PageOffsetCalculation(req.Page, req.PageSize)
	users, total, err := models.ListUser(s.server.Gorm(ctx.Context()), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("user list: %w", err)
	}
	list := make([]entity.User, 0, len(users))
	for k := range users {
		list = append(list, *entity.NewUserFromModel(&users[k], "", ""))
	}
	resp := &entity.UserListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}
	return resp, nil
}
