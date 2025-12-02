//go:build wireinject
// +build wireinject

package wire

import (
	"es_test/Service"
	"es_test/common"
	"es_test/conf"
	"es_test/dao"
	"es_test/handler"

	"github.com/google/wire"
)

func InitializeHandler(conf *conf.Config) *handler.UserHandler {
	wire.Build(
		common.NewEsClient,
		common.NewRouterClient,

		dao.NewUserES,
		Service.NewUserService,

		handler.NewUserHandler)
	return &handler.UserHandler{}
}
