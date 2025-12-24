package service

import (
	"context"

	"github.com/HoronLee/EchoHub/internal/model/helloworld"
)

// HelloWorldRepo 定义HelloWorld数据访问接口
type HelloWorldRepo interface {
	CreateHelloWorld(ctx context.Context, hw *helloworld.HelloWorld) error
	GetDatabaseInfo(ctx context.Context) (string, error)
}

// HelloWorldService HelloWorld服务实现
type HelloWorldService struct {
	repo HelloWorldRepo
}

// NewHelloWorldService 创建HelloWorldService实例
func NewHelloWorldService(repo HelloWorldRepo) *HelloWorldService {
	return &HelloWorldService{repo: repo}
}

func (s *HelloWorldService) PostHelloWorld(ctx context.Context, message string) error {
	hw := &helloworld.HelloWorld{Message: message}
	return s.repo.CreateHelloWorld(ctx, hw)
}

func (s *HelloWorldService) GetDatabaseInfo(ctx context.Context) (string, error) {
	return s.repo.GetDatabaseInfo(ctx)
}
