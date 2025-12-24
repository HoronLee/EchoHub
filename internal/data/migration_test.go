package data

import (
	"os"
	"testing"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/user"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/stretchr/testify/assert"
)

func TestUserModelMigration(t *testing.T) {
	// 创建测试配置，使用内存 SQLite 数据库
	cfg := &config.AppConfig{}
	cfg.Database.Driver = "sqlite"
	cfg.Database.Source = ":memory:"
	cfg.Server.Mode = "debug"

	// 创建测试日志器
	logger := util.NewLogger(cfg)

	// 初始化数据库（包含 AutoMigrate）
	db, err := NewDB(cfg, logger)
	assert.NoError(t, err, "Database initialization should succeed")
	assert.NotNil(t, db, "Database instance should not be nil")

	// 验证 users 表是否存在
	hasTable := db.Migrator().HasTable(&user.User{})
	assert.True(t, hasTable, "users table should exist after migration")

	// 验证表结构是否正确
	hasColumn := db.Migrator().HasColumn(&user.User{}, "id")
	assert.True(t, hasColumn, "users table should have id column")

	hasColumn = db.Migrator().HasColumn(&user.User{}, "username")
	assert.True(t, hasColumn, "users table should have username column")

	hasColumn = db.Migrator().HasColumn(&user.User{}, "password")
	assert.True(t, hasColumn, "users table should have password column")

	hasColumn = db.Migrator().HasColumn(&user.User{}, "created_at")
	assert.True(t, hasColumn, "users table should have created_at column")

	hasColumn = db.Migrator().HasColumn(&user.User{}, "updated_at")
	assert.True(t, hasColumn, "users table should have updated_at column")

	// 验证索引是否存在（GORM 自动创建的索引名为 idx_users_username）
	hasIndex := db.Migrator().HasIndex(&user.User{}, "idx_users_username")
	assert.True(t, hasIndex, "users table should have unique index on username")

	// 测试基本的 CRUD 操作以确保表结构正确
	testUser := &user.User{
		Username: "testuser",
		Password: "hashedpassword",
	}

	// 创建用户
	result := db.Create(testUser)
	assert.NoError(t, result.Error, "Should be able to create user")
	assert.NotZero(t, testUser.ID, "User ID should be set after creation")

	// 查询用户
	var foundUser user.User
	result = db.Where("username = ?", "testuser").First(&foundUser)
	assert.NoError(t, result.Error, "Should be able to find user")
	assert.Equal(t, "testuser", foundUser.Username, "Username should match")
	assert.Equal(t, "hashedpassword", foundUser.Password, "Password should match")

	// 测试唯一约束
	duplicateUser := &user.User{
		Username: "testuser", // 相同的用户名
		Password: "anotherpassword",
	}
	result = db.Create(duplicateUser)
	assert.Error(t, result.Error, "Should not be able to create user with duplicate username")

	// 清理
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func TestDatabaseMigrationWithExistingData(t *testing.T) {
	// 使用临时文件而不是内存数据库来测试数据持久性
	tempFile := "/tmp/test_migration.db"
	defer os.Remove(tempFile) // 清理临时文件

	// 创建测试配置
	cfg := &config.AppConfig{}
	cfg.Database.Driver = "sqlite"
	cfg.Database.Source = tempFile
	cfg.Server.Mode = "debug"

	logger := util.NewLogger(cfg)

	// 第一次初始化数据库
	db1, err := NewDB(cfg, logger)
	assert.NoError(t, err)

	// 创建一些测试数据
	testUser := &user.User{
		Username: "existinguser",
		Password: "hashedpassword",
	}
	result := db1.Create(testUser)
	assert.NoError(t, result.Error)

	// 关闭第一个连接
	sqlDB1, _ := db1.DB()
	sqlDB1.Close()

	// 重新初始化数据库（模拟应用重启）
	db2, err := NewDB(cfg, logger)
	assert.NoError(t, err)

	// 验证数据仍然存在且表结构正确
	var foundUser user.User
	result = db2.Where("username = ?", "existinguser").First(&foundUser)
	assert.NoError(t, result.Error, "Existing data should be preserved after migration")
	assert.Equal(t, "existinguser", foundUser.Username)

	// 清理
	sqlDB2, _ := db2.DB()
	sqlDB2.Close()
}
