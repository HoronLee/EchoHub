package util

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger Gorm日志适配器
type GormLogger struct {
	logger                    *Logger
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建Gorm日志适配器
func NewGormLogger(logger *Logger) *GormLogger {
	return &GormLogger{
		logger:                    logger,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		fields = append(fields, zap.Error(err))
		l.logger.Error("Database Error", fields...)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.Warn("Slow SQL", fields...)
		return
	}

	l.logger.Debug("Database Query", fields...)
}
