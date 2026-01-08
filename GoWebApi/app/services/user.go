package services

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/orm/orm"
	"gorm.io/gorm"

	"webapi/app/entity"
	"webapi/app/models"
	"webapi/app/models/dao"
	"webapi/core"
)

// User 用户相关
type User struct {
	server *core.Server
	dao    *dao.Dao
}

func NewUser(server *core.Server, modelDao *dao.Dao) *User {
	return &User{
		server: server,
		dao:    modelDao,
	}
}

// Update 更新用户信息
func (s *User) Update(ctx fiber.Ctx, req *entity.UserUpdateRequest) (*entity.User, error) {
	user, err := s.dao.GetUserByID(ctx.Context(), req.Uid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
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
	err = s.dao.UpdateUser(ctx.Context(), user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Info 获取用户信息
func (s *User) Info(ctx fiber.Ctx, req *entity.UserInfoRequest) (*entity.User, error) {
	user, err := s.dao.GetUserByID(ctx.Context(), req.Uid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "用户不存在")
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
	limit, offset := orm.PageLimitOffset(req.Page, req.PageSize)
	users, total, err := s.dao.ListUser(ctx.Context(), limit, offset)
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
