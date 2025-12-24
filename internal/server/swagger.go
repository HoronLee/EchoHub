package server

import (
	"fmt"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/swagger"
	envUtil "github.com/HoronLee/EchoHub/internal/util/env"
)

// configureSwagger 根据配置和环境变量动态配置Swagger信息
func configureSwagger(cfg *config.AppConfig) {
	// 从环境变量或配置文件获取Swagger设置
	host := envUtil.GetEnvOrConfig("SWAGGER_HOST", cfg.Swagger.Host)
	if host == "" {
		// 如果没有配置，使用服务器配置构建默认host
		host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		if cfg.Server.Host == "0.0.0.0" {
			host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
		}
	}

	basePath := envUtil.GetEnvOrConfig("SWAGGER_BASE_PATH", cfg.Swagger.BasePath)
	if basePath == "" {
		basePath = "/api"
	}

	schemes := cfg.Swagger.Schemes
	if len(schemes) == 0 {
		if cfg.Server.Mode == "release" {
			schemes = []string{"https", "http"}
		} else {
			schemes = []string{"http", "https"}
		}
	}

	title := envUtil.GetEnvOrConfig("SWAGGER_TITLE", cfg.Swagger.Title)
	if title == "" {
		title = "EchoHub API 文档"
	}

	description := envUtil.GetEnvOrConfig("SWAGGER_DESCRIPTION", cfg.Swagger.Description)
	if description == "" {
		description = "基于Gin、Gorm、Viper、Wire、Cobra的HTTP快速开发框架 API 文档"
	}

	version := envUtil.GetEnvOrConfig("SWAGGER_VERSION", cfg.Swagger.Version)
	if version == "" {
		version = "1.0"
	}

	// 更新Swagger文档信息
	swagger.SwaggerInfo.Host = host
	swagger.SwaggerInfo.BasePath = basePath
	swagger.SwaggerInfo.Schemes = schemes
	swagger.SwaggerInfo.Title = title
	swagger.SwaggerInfo.Description = description
	swagger.SwaggerInfo.Version = version
}
