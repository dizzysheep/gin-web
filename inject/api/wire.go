//go:build wireinject
// +build wireinject

package api

import (
	"gin-web/app/handler"
	"gin-web/internal/dao"
	"gin-web/internal/service"
	"gin-web/pkg/boostrap"
	"github.com/google/wire"
)

type AppContainer struct {
	Handlers *handler.Handlers
}

func NewAppContainer() *AppContainer {
	wire.Build(
		boostrap.InitDBEngine,

		//实例化Daos
		dao.ProviderSet,
		wire.Struct(new(dao.Daos), "*"),

		//实例化Services
		service.ProviderSet,
		wire.Struct(new(service.Services), "*"),

		//实例化Handlers
		handler.ProviderSet,
		wire.Struct(new(handler.Handlers), "*"),

		wire.Struct(new(AppContainer), "*"),
	)
	return nil
}
