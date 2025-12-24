package env

import "os"

// GetEnvOrConfig 优先从环境变量获取值，如果不存在则使用配置文件的值
func GetEnvOrConfig(envKey, configValue string) string {
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return configValue
}
