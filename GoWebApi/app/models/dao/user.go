package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"webapi/app/models"
)

func (d *Dao) CreateUser(ctx context.Context, user *models.User) error {
	// 事务（返回任何错误都会回滚事务）
	return d.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userExtend := &models.UserExtend{
			Uid: user.ID,
		}
		if err := tx.Create(userExtend).Error; err != nil {
			return err
		}
		userThirdParty := &models.UserThirdParty{
			Uid: user.ID,
		}
		if err := tx.Create(userThirdParty).Error; err != nil {
			return err
		}
		// INSERT INTO `user` (`type_id`,`mobile`,`email`,`password`,`nickname`,`avatar_url`,`states`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)
		// INSERT INTO `user_extend` (`uid`,`birthday`,`height`,`weight`) VALUES (?,?,?)
		// INSERT INTO `user_third_party` (`uid`,`wechat_unionid`,`wechat_openid`) VALUES (?,?,?)
		return nil
	})
}

func (d *Dao) UpdateUser(ctx context.Context, user *models.User) error {
	updates := make(map[string]any)
	if user.TypeId > 0 {
		updates["type_id"] = user.TypeId
	}
	if user.Mobile != "" {
		updates["mobile"] = user.Mobile
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Nickname != "" {
		updates["nickname"] = user.Nickname
	}
	if user.AvatarUrl != "" {
		updates["avatar_url"] = user.AvatarUrl
	}
	if user.States > 0 {
		updates["states"] = user.States
	}
	// UPDATE `user` SET `type_id`=?,`mobile`=?,`email`=?,`password`=?,`nickname`=?,`avatar_url`=?,`states`=?,`updated_at`=? WHERE `id` = ?
	return d.UpdatesById(ctx, user.ID, d.User, updates)
}

func (d *Dao) GetUser(ctx context.Context, uid int, email, mobile, nickname string) (*models.User, error) {
	if uid <= 0 && email == "" && mobile == "" && nickname == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var user models.User
	client := d.WithContext(ctx)
	if uid > 0 {
		client = client.Where("id", uid)
	}
	if email != "" {
		client = client.Where("email", email).Order("id")
	}
	if mobile != "" {
		client = client.Where("mobile", mobile).Order("id")
	}
	if nickname != "" {
		client = client.Where("nickname", nickname).Order("id")
	}
	if err := client.Take(&user).Error; err != nil {
		return nil, err
	}
	// SELECT * FROM `user` WHERE `id` = ? LIMIT 1
	// SELECT * FROM `user` WHERE `email` = ? ORDER BY id LIMIT 1
	// SELECT * FROM `user` WHERE `mobile` = ? ORDER BY id LIMIT 1
	// SELECT * FROM `user` WHERE `nickname` = ? ORDER BY id LIMIT 1
	return &user, nil
}

func (d *Dao) GetUserByID(ctx context.Context, uid int) (*models.User, error) {
	return d.GetUser(ctx, uid, "", "", "")
}

func (d *Dao) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return d.GetUser(ctx, 0, email, "", "")
}

func (d *Dao) GetUserByMobile(ctx context.Context, mobile string) (*models.User, error) {
	return d.GetUser(ctx, 0, "", mobile, "")
}

func (d *Dao) GetUserByNickname(ctx context.Context, nickname string) (*models.User, error) {
	return d.GetUser(ctx, 0, "", "", nickname)
}

func (d *Dao) ListUser(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
	var total int64
	var users []models.User
	db := d.WithContext(ctx)
	if err := db.Model(d.User).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total <= 0 {
		return users, total, nil
	}
	if err := db.Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true}).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	// SELECT count(*) FROM `user`
	// SELECT * FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
	return users, total, nil
}
