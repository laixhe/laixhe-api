package services

import (
	"github.com/gofiber/fiber/v2"

	"webapi/app/entity"
	"webapi/app/model/dao"
	"webapi/core"
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

// Login 登录
func (a *Auth) Login(c *fiber.Ctx, req *entity.AuthLoginRequest) (*entity.AuthLoginResponse, error) {
	user, err := a.dao.GetUserByNickname(c.UserContext(), "laixhe")
	if err != nil {
		return nil, err
	}
	return &entity.AuthLoginResponse{
		Token: user.Nickname,
	}, nil
}
