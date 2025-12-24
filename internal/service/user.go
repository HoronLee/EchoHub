package service

import (
	"context"
	"errors"
	"time"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/user"
	cryptoUtil "github.com/HoronLee/EchoHub/internal/util/crypto"
	jwtutil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// UserRepo 定义用户数据访问接口
type UserRepo interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByID(ctx context.Context, id uint) (*user.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

// UserService 用户服务实现
type UserService struct {
	repo      UserRepo
	jwtHelper *jwtutil.JWT[user.Claims]
}

// NewUserService 创建UserService实例（通过Wire注入）
func NewUserService(repo UserRepo) *UserService {
	// 创建JWT helper
	jwtCfg := &jwtutil.Config{
		SecretKey: string(config.JWT_SECRET),
	}
	jwtHelper := jwtutil.NewJWT[user.Claims](jwtCfg)

	return &UserService{
		repo:      repo,
		jwtHelper: jwtHelper,
	}
}

// Register 用户注册
// 检查用户名是否存在，使用MD5加密密码，创建用户
func (s *UserService) Register(ctx context.Context, req user.RegisterRequest) error {
	// 1. 检查用户名是否已存在
	existingUser, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库查询错误
		return err
	}
	if existingUser != nil {
		// 用户名已存在
		return errors.New("username already exists")
	}

	// 2. 使用MD5加密密码
	hashedPassword := cryptoUtil.MD5Encrypt(req.Password)

	// 3. 创建用户
	newUser := &user.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	return s.repo.CreateUser(ctx, newUser)
}

// Login 用户登录
// 验证用户名和密码，生成JWT Token
func (s *UserService) Login(ctx context.Context, req user.LoginRequest) (string, error) {
	// 1. 查询用户
	u, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// 2. 验证密码
	hashedPassword := cryptoUtil.MD5Encrypt(req.Password)
	if u.Password != hashedPassword {
		return "", errors.New("invalid username or password")
	}

	// 3. 生成JWT Token
	claims := &user.Claims{
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Config.Auth.Jwt.Expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-60 * time.Second)), // 允许60秒的时钟偏差
			Issuer:    config.Config.Auth.Jwt.Issuer,
			Subject:   u.Username,
			Audience:  []string{config.Config.Auth.Jwt.Audience},
		},
	}

	token, err := s.jwtHelper.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// 1. 检查用户是否存在
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// 2. 删除用户
	return s.repo.DeleteUser(ctx, userID)
}
