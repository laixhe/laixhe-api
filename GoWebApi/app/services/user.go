package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/orm/orm"
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

func NewUser(server *core.Server) *User {
	return &User{
		server: server,
	}
}

// Update 更新用户信息
func (s *User) Update(ctx fiber.Ctx, req *entity.UserUpdateRequest) (*entity.User, error) {
	user := &models.User{}
	if err := s.server.Orm().GetById(ctx.Context(), req.Uid, user); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.AuthorizedError()
	}
	resp := &entity.User{
		Uid:       user.ID,
		TypeId:    user.TypeId,
		Account:   user.Account,
		Mobile:    user.Mobile,
		Email:     user.Email,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
		Sex:       user.Sex,
		States:    user.States,
		CreatedAt: user.CreatedAt.Format(time.DateTime),
	}
	user = &models.User{
		ID:        user.ID,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
	}
	if err := user.Update(s.server.Gorm(ctx.Context())); err != nil {
		return nil, err
	}
	return resp, nil
}

// Info 获取用户信息
func (s *User) Info(ctx fiber.Ctx, req *entity.UserInfoRequest) (*entity.User, error) {
	user := &models.User{}
	if err := s.server.Orm().GetById(ctx.Context(), req.Uid, user); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, xfiber.ParamError("用户不存在")
	}
	return &entity.User{
		Uid:       user.ID,
		TypeId:    user.TypeId,
		Account:   user.Account,
		Mobile:    user.Mobile,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarUrl: user.AvatarUrl,
		Sex:       user.Sex,
		States:    user.States,
		CreatedAt: user.CreatedAt.Format(time.DateTime),
	}, nil
}

// List 获取用户列表
func (s *User) List(ctx fiber.Ctx, req *entity.UserListRequest) (*entity.UserListResponse, error) {
	limit, offset := orm.PageOffsetCalculation(req.Page, req.PageSize)
	users, total, err := new(models.User).List(s.server.Gorm(ctx.Context()), limit, offset)
	if err != nil {
		return nil, err
	}
	list := make([]entity.User, 0, len(users))
	for k := range users {
		list = append(list, entity.User{
			Uid:       users[k].ID,
			TypeId:    users[k].TypeId,
			Account:   users[k].Account,
			Mobile:    users[k].Mobile,
			Email:     users[k].Email,
			Nickname:  users[k].Nickname,
			AvatarUrl: users[k].AvatarUrl,
			Sex:       users[k].Sex,
			States:    users[k].States,
			CreatedAt: users[k].CreatedAt.Format(time.DateTime),
		})
	}
	resp := &entity.UserListResponse{
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}
	return resp, nil
}
