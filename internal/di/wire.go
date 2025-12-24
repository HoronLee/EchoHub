//go:build wireinject
// +build wireinject

package di

import (
	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/data"
	"github.com/HoronLee/EchoHub/internal/handler"
	"github.com/HoronLee/EchoHub/internal/server"
	"github.com/HoronLee/EchoHub/internal/service"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/google/wire"
)

// InitServer 初始化服务器
func InitServer(cfg *config.AppConfig) (*server.HTTPServer, func(), error) {
	wire.Build(
		util.NewLogger,
		data.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		server.ProviderSet,
	)
	return nil, nil, nil
}
