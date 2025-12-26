// Package server
//
// @title EchoHub API 文档
// @version 1.0
// @description 基于Echo、Gorm、Viper、Wire、Cobra的HTTP快速开发框架 API 文档
// @contact.name EchoHub Team
// @contact.url https://github.com/HoronLee/EchoHub
// @contact.email support@echohub.dev
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
// @schemes http https
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/handler"
	"github.com/HoronLee/EchoHub/internal/middleware"
	"github.com/HoronLee/EchoHub/internal/router"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/HoronLee/EchoHub/internal/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HTTPServer struct {
	cfg        *config.AppConfig
	echo       *echo.Echo
	httpServer *http.Server
	handlers   *handler.Handlers
	db         *gorm.DB
	logger     *util.Logger
	validator  *validator.Validator
}

func NewHTTPServer(
	cfg *config.AppConfig,
	handlers *handler.Handlers,
	db *gorm.DB,
	logger *util.Logger,
	v *validator.Validator,
) *HTTPServer {
	e := echo.New()

	if cfg.Server.Mode == "release" {
		e.HideBanner = true
		e.Debug = false
	} else {
		e.Debug = true
	}

	// 配置Swagger信息
	configureSwagger(cfg)

	// 配置验证器
	v.SetupEcho(e)

	// 设置自定义错误处理器（必须在中间件之前）
	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler

	// 中间件
	e.Use(middleware.Logger(logger))
	e.Use(middleware.Recovery(logger))
	e.Use(middleware.CORS(cfg))

	return &HTTPServer{
		cfg:       cfg,
		echo:      e,
		handlers:  handlers,
		db:        db,
		logger:    logger,
		validator: v,
	}
}

func (s *HTTPServer) Start() error {
	router.SetupRouter(s.echo, s.handlers)

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.echo,
	}

	s.logger.Info("Server starting", zap.String("addr", addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *HTTPServer) GetEcho() *echo.Echo {
	return s.echo
}
