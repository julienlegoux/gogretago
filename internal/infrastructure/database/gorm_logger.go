package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lgxju/gogretago/internal/lib/shared"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const slowQueryThreshold = 200 * time.Millisecond

// gormLogAdapter bridges GORM's logger interface to the application's structured
// logger, avoiding ANSI color codes in log output.
type gormLogAdapter struct {
	logger   *shared.Logger
	logLevel gormlogger.LogLevel
}

func newGormLogger(appLogger *shared.Logger, level gormlogger.LogLevel) gormlogger.Interface {
	return &gormLogAdapter{
		logger:   appLogger,
		logLevel: level,
	}
}

func (l *gormLogAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &gormLogAdapter{
		logger:   l.logger,
		logLevel: level,
	}
}

func (l *gormLogAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...), nil)
	}
}

func (l *gormLogAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...), nil)
	}
}

func (l *gormLogAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...), nil)
	}
}

func (l *gormLogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logCtx := map[string]interface{}{
		"elapsed_ms": float64(elapsed.Nanoseconds()) / 1e6,
		"rows":       rows,
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.logLevel >= gormlogger.Error:
		logCtx["error"] = err.Error()
		logCtx["sql"] = sql
		l.logger.Error("gorm query error", logCtx)
	case elapsed > slowQueryThreshold && l.logLevel >= gormlogger.Warn:
		logCtx["sql"] = sql
		l.logger.Warn("gorm slow query", logCtx)
	case l.logLevel >= gormlogger.Info:
		logCtx["sql"] = sql
		l.logger.Debug("gorm query", logCtx)
	}
}
