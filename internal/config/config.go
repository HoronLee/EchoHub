package config

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"log"
	"os"

	"github.com/spf13/viper"

	model "github.com/HoronLee/EchoHub/internal/model/common"
)

// Config 全局配置变量
var Config AppConfig

// JWT_SECRET 用于JWT签名的密钥
var JWT_SECRET []byte

// AppConfig 应用程序配置结构体
type AppConfig struct {
	Server struct {
		Port string `mapstructure:"port"` // 服务器端口
		Host string `mapstructure:"host"` // 服务器主机地址
		Mode string `mapstructure:"mode"` // 运行模式，可能的值为 "debug" 或 "release"
	} `mapstructure:"server"`
	Database struct {
		Driver string `mapstructure:"type"`   // 数据库驱动
		Source string `mapstructure:"source"` // 数据库连接字符串
	} `mapstructure:"database"`
	Auth struct {
		Jwt struct {
			Secret   string `mapstructure:"secret"`   // JWT的密钥
			Expires  int    `mapstructure:"expires"`  // JWT的过期时间，单位为秒
			Issuer   string `mapstructure:"issuer"`   // JWT的发行者
			Audience string `mapstructure:"audience"` // JWT的受众
		} `mapstructure:"jwt"`
	} `mapstructure:"auth"`
	Swagger struct {
		Host         string   `mapstructure:"host"`          // Swagger文档的主机地址
		BasePath     string   `mapstructure:"basepath"`      // API基础路径
		Schemes      []string `mapstructure:"schemes"`       // 支持的协议方案
		Title        string   `mapstructure:"title"`         // API文档标题
		Description  string   `mapstructure:"description"`   // API文档描述
		Version      string   `mapstructure:"version"`       // API版本
		ContactName  string   `mapstructure:"contact_name"`  // 联系人姓名
		ContactURL   string   `mapstructure:"contact_url"`   // 联系人URL
		ContactEmail string   `mapstructure:"contact_email"` // 联系人邮箱
		LicenseName  string   `mapstructure:"license_name"`  // 许可证名称
		LicenseURL   string   `mapstructure:"license_url"`   // 许可证URL
	} `mapstructure:"swagger"`
	CORS struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的来源列表
		AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的HTTP方法
		AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的请求头
		ExposeHeaders    []string `mapstructure:"expose_headers"`    // 暴露的响应头
		AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许发送凭证
		MaxAge           int      `mapstructure:"max_age"`           // 预检请求缓存时间（秒）
	} `mapstructure:"cors"`
}

//go:embed config.yaml
var configData []byte

// LoadAppConfig 加载应用程序配置
// configPath: 外部配置文件路径，如果为空则只使用嵌入式配置
func LoadAppConfig(configPath string) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 1. 先加载嵌入式配置作为默认配置
	err := v.ReadConfig(bytes.NewReader(configData))
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	// 2. 如果指定了外部配置文件，则用外部配置覆盖
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			v.SetConfigFile(configPath)
			err = v.MergeInConfig() // 使用MergeInConfig合并配置
			if err != nil {
				log.Printf("Warning: failed to merge external config: %v, using embedded config\n", err)
			} else {
				log.Printf("Loaded external config from: %s\n", configPath)
			}
		} else {
			log.Printf("Warning: external config file not found: %s, using embedded config\n", configPath)
		}
	}

	// 3. 将配置反序列化到结构体
	err = v.Unmarshal(&Config)
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	// 4. 初始化 JWT_SECRET
	JWT_SECRET = GetJWTSecret()
}

// GetJWTSecret 加载JWT密钥
func GetJWTSecret() []byte {
	// 优先级：环境变量 > 配置文件 > 随机生成
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 从配置文件中获取
		secret = Config.Auth.Jwt.Secret
	}
	if secret == "" {
		// 如果都没有设置，则使用随机生成的密钥
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal("failed to generate random JWT secret:", err)
		}
		secret = hex.EncodeToString(b)
		log.Println("Warning: Using randomly generated JWT secret. Set JWT_SECRET environment variable or auth.jwt.secret in config for production.")
	}

	return []byte(secret)
}
