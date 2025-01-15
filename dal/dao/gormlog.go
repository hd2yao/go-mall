package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/hd2yao/go-mall/common/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration // 慢 SQL 阈值，也可以做成配置项放到配置文件中
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: time.Millisecond * 500,
	}
}

// 实现 gormLogger.Interface 接口

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.New(ctx).Info(msg, "data", data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.New(ctx).Warn(msg, "data", data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.New(ctx).Error(msg, "data", data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	duration := time.Since(begin).Milliseconds()
	// 获取 SQL 语句和返回条数
	sql, rows := fc()

	// Gorm 错误时记录错误日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.New(ctx).Error("SQL ERROR", "sql", sql, "rows", rows, "dur(ms)", duration)
	}

	// 慢 SQL 日志
	if duration > l.SlowThreshold.Milliseconds() {
		logger.New(ctx).Warn("SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", duration)
	} else {
		logger.New(ctx).Debug("SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
}
