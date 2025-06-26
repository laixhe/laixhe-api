package services

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"webapi/app/entity"
	"webapi/app/model"
	"webapi/app/model/dao"
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
func (s *User) Update(ctx *fiber.Ctx, req *entity.UserUpdateRequest) (*entity.User, error) {
	user, err := s.dao.GetUserByID(ctx.UserContext(), req.Uid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	resp := &entity.User{
		Uid:       user.ID,
		TypeId:    user.TypeId,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
		States:    user.States,
	}
	user = &model.User{
		ID:        user.ID,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
	}
	err = s.dao.UpdateUser(ctx.UserContext(), user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Info 获取用户信息
func (s *User) Info(ctx *fiber.Ctx, req *entity.UserInfoRequest) (*entity.User, error) {
	user, err := s.dao.GetUserByID(ctx.UserContext(), req.Uid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "用户不存在")
	}
	return &entity.User{
		Uid:       user.ID,
		TypeId:    user.TypeId,
		Nickname:  user.Nickname,
		AvatarUrl: user.AvatarUrl,
		States:    user.States,
	}, nil
}

// List 获取用户列表
func (s *User) List(ctx *fiber.Ctx, req *entity.UserListRequest) (*entity.UserListResponse, error) {
	count := s.dao.ListUserCount(ctx.UserContext())
	users, err := s.dao.ListUser(ctx.UserContext(), req.OffsetId, req.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]entity.User, 0, len(users))
	for k := range users {
		list = append(list, entity.User{
			Uid:       users[k].ID,
			TypeId:    users[k].TypeId,
			Nickname:  users[k].Nickname,
			AvatarUrl: users[k].AvatarUrl,
			States:    users[k].States,
		})
	}
	resp := &entity.UserListResponse{
		Total: int(count),
		List:  list,
	}
	return resp, nil
}
