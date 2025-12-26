package data

import (
	"context"

	"github.com/HoronLee/EchoHub/internal/model/user"
	"github.com/HoronLee/EchoHub/internal/service"
	"github.com/HoronLee/EchoHub/internal/util/log"
	"go.uber.org/zap"
)

// userRepo 用户数据访问实现
type userRepo struct {
	data *Data
	log  *log.Logger
}

// NewUserRepo 创建UserRepo实例
// 注意：返回的是 service.UserRepo 接口类型
func NewUserRepo(data *Data, logger *log.Logger) service.UserRepo {
	return &userRepo{
		data: data,
		log:  logger,
	}
}

// CreateUser 创建用户记录
func (r *userRepo) CreateUser(ctx context.Context, u *user.User) error {
	r.log.Debug("Creating user", zap.String("username", u.Username))
	err := r.data.db.WithContext(ctx).Create(u).Error
	if err != nil {
		r.log.Error("Failed to create user", zap.Error(err), zap.String("username", u.Username))
		return err
	}
	r.log.Info("User created successfully", zap.String("username", u.Username), zap.Uint("id", u.ID))
	return nil
}

// GetUserByUsername 根据用户名查询用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	r.log.Debug("Getting user by username", zap.String("username", username))
	var u user.User
	err := r.data.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		r.log.Debug("User not found", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	r.log.Debug("User found", zap.String("username", username), zap.Uint("id", u.ID))
	return &u, nil
}

// GetUserByID 根据用户ID查询用户
func (r *userRepo) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	r.log.Debug("Getting user by ID", zap.Uint("id", id))
	var u user.User
	err := r.data.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		r.log.Debug("User not found", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	r.log.Debug("User found", zap.Uint("id", id), zap.String("username", u.Username))
	return &u, nil
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id uint) error {
	r.log.Debug("Deleting user", zap.Uint("id", id))
	err := r.data.db.WithContext(ctx).Delete(&user.User{}, id).Error
	if err != nil {
		r.log.Error("Failed to delete user", zap.Error(err), zap.Uint("id", id))
		return err
	}
	r.log.Info("User deleted successfully", zap.Uint("id", id))
	return nil
}
