package handler

import "github.com/google/wire"

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(NewHandlers, NewHelloWorldHandler, NewUserHandler)

// Handlers 聚合各个模块的Handler
type Handlers struct {
	HelloWorldHandler *HelloWorldHandler
	UserHandler       *UserHandler
}

// NewHandlers 创建Handlers实例
func NewHandlers(hwHandler *HelloWorldHandler, userHandler *UserHandler) *Handlers {
	return &Handlers{
		HelloWorldHandler: hwHandler,
		UserHandler:       userHandler,
	}
}
