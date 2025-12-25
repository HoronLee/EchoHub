package model

var (
	// Version 版本号，构建时通过 -ldflags 注入
	Version = "0.1.0"
	// BuildTime 构建时间，构建时通过 -ldflags 注入
	BuildTime = "unknown"
	// GitCommit Git 提交哈希，构建时通过 -ldflags 注入
	GitCommit = "unknown"
)
