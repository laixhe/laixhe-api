package services

import (
	"webapi/app/models"
)

type Service struct {
	data *models.Data
}

func NewService() *Service {
	return &Service{
		data: models.NewData(),
	}
}
