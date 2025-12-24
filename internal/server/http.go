// Package server
//
// @title EchoHub API 文档
// @version 1.0
// @description 基于Gin、Gorm、Viper、Wire、Cobra的HTTP快速开发框架 API 文档
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
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HTTPServer struct {
	cfg        *config.AppConfig
	engine     *gin.Engine
	httpServer *http.Server
	handlers   *handler.Handlers
	db         *gorm.DB
	logger     *util.Logger
}

func NewHTTPServer(
	cfg *config.AppConfig,
	handlers *handler.Handlers,
	db *gorm.DB,
	logger *util.Logger,
) *HTTPServer {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 配置Swagger信息
	configureSwagger(cfg)

	engine := gin.New()
	engine.Use(middleware.Logger(logger))
	engine.Use(middleware.Recovery(logger))

	return &HTTPServer{
		cfg:      cfg,
		engine:   engine,
		handlers: handlers,
		db:       db,
		logger:   logger,
	}
}

func (s *HTTPServer) Start() error {
	router.SetupRouter(s.engine, s.handlers)

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
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

func (s *HTTPServer) GetEngine() *gin.Engine {
	return s.engine
}
