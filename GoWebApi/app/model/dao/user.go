package dao

import (
	"context"

	"webapi/app/model"
)

func (d *Dao) GetUserByNickname(ctx context.Context, nickname string) (*model.User, error) {
	var user model.User
	if err := d.WithContext(ctx).
		Where("nickname", nickname).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
