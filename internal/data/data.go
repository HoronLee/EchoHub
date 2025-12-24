package data

import (
	"fmt"
	"time"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/helloworld"
	"github.com/HoronLee/EchoHub/internal/model/user"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewHelloWorldRepo, NewUserRepo)

// Data 统一的数据访问层结构体
type Data struct {
	db  *gorm.DB
	log *util.Logger
}

// NewData 创建Data实例
func NewData(db *gorm.DB, logger *util.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Info("closing the data resources")
	}
	return &Data{
		db:  db,
		log: logger,
	}, cleanup, nil
}

// NewDB 创建数据库连接
func NewDB(cfg *config.AppConfig, logger *util.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// 根据配置选择数据库驱动
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Database.Source)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.Source)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// 配置GORM日志
	gormLogger := util.NewGormLogger(logger)

	// 打开数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connected successfully", zap.String("driver", cfg.Database.Driver))

	// 自动迁移数据库表
	if err = db.AutoMigrate(
		&helloworld.HelloWorld{},
		&user.User{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
