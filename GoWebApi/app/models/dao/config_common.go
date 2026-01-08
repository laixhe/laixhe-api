package dao

import (
	"context"
	
	"webapi/app/models"
)

func (d *Dao) ListConfigCommon(ctx context.Context, keys ...string) ([]models.ConfigCommon, error) {
	var list []models.ConfigCommon
	client := d.WithContext(ctx)
	if len(keys) > 0 {
		client = client.Where("key IN ?", keys)
	}
	err := client.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
