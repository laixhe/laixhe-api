package services

import (
	"errors"

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
func (s *User) Update(ctx fiber.Ctx, req *entity.UserUpdateRequest) (*entity.User, error) {
	user := &models.User{}
	if err := s.server.Orm().GetById(ctx.Context(), user, req.Uid); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.ParamError("用户不存在")
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
		return nil, err
	}
	return resp, nil
}

// Info 获取用户信息
func (s *User) Info(ctx fiber.Ctx, req *entity.UserInfoRequest) (*entity.User, error) {
	user := &models.User{}
	if err := s.server.Orm().GetById(ctx.Context(), user, req.Uid); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.ParamError("用户不存在")
	}
	return entity.NewUserFromModel(user, "", ""), nil
}

// List 获取用户列表
func (s *User) List(ctx fiber.Ctx, req *entity.UserListRequest) (*entity.UserListResponse, error) {
	limit, offset := orm.PageOffsetCalculation(req.Page, req.PageSize)
	users, total, err := models.ListUser(s.server.Gorm(ctx.Context()), limit, offset)
	if err != nil {
		return nil, err
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
