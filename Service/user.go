package Service

import (
	"context"
	"es_test/dao"
	"es_test/model"
)

type UserService struct {
	es *dao.UserES
}

func NewUserService(es *dao.UserES) *UserService {
	return &UserService{
		es: es,
	}
}

func (s *UserService) BatchAdd(ctx context.Context, user []*model.UserEs) error {
	return s.es.BatchAdd(ctx, user)
}
