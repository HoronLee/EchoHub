package data

import (
	"context"

	"github.com/HoronLee/EchoHub/internal/model/helloworld"
	"github.com/HoronLee/EchoHub/internal/service"
	"github.com/HoronLee/EchoHub/internal/util/log"
	"go.uber.org/zap"
)

// helloworldRepo HelloWorld数据访问实现
type helloworldRepo struct {
	data *Data
	log  *log.Logger
}

// NewHelloWorldRepo 创建HelloWorldRepo实例
func NewHelloWorldRepo(data *Data, logger *log.Logger) service.HelloWorldRepo {
	return &helloworldRepo{
		data: data,
		log:  logger,
	}
}

// CreateHelloWorld 创建HelloWorld记录
func (r *helloworldRepo) CreateHelloWorld(ctx context.Context, hw *helloworld.HelloWorld) error {
	r.data.log.Debug("Creating HelloWorld record", zap.String("message", hw.Message))
	err := r.data.db.WithContext(ctx).Create(hw).Error
	if err != nil {
		r.data.log.Error("Failed to create HelloWorld record", zap.Error(err))
		return err
	}
	r.data.log.Info("HelloWorld record created successfully", zap.Uint("id", hw.ID))
	return nil
}

// GetDatabaseInfo 获取数据库连接信息
func (r *helloworldRepo) GetDatabaseInfo(ctx context.Context) (string, error) {
	r.data.log.Debug("Getting database info")
	dbName := r.data.db.Migrator().CurrentDatabase()
	if dbName == "" {
		dbName = "unknown"
	}
	info := "Connected to: " + dbName
	r.data.log.Debug("Database info retrieved", zap.String("info", info))
	return info, nil
}
