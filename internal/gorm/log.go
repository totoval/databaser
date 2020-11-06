package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/totoval/framework/helpers/toto"
	facade_logger "github.com/totoval/logger/pkg/facade"
	"gorm.io/gorm/logger"
	"time"
)

type Log struct {
	slowThreshold time.Duration
	logger        facade_logger.Logger
}

func NewLog(logger facade_logger.Logger, slowThreshold time.Duration) *Log {
	return &Log{slowThreshold: slowThreshold, logger: logger}
}

func (l *Log) SlowThreshold() time.Duration {
	if l.slowThreshold == 0 {
		return 200 * time.Millisecond
	}
	return l.slowThreshold
}

func (l *Log) LogMode(level logger.LogLevel) logger.Interface {
	// log level will use outter logger level, this level will be useless
	return l
}

func (l *Log) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Info(s, toto.V{"data": i})
}

func (l *Log) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Warn(s, toto.V{"data": i})
}

func (l *Log) Error(ctx context.Context, s string, i ...interface{}) {
	_ = l.logger.Error(errors.New(s), toto.V{"data": i})
}

func (l *Log) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	msg := ""
	switch {
	case err != nil:
		msg = err.Error()
	case elapsed > l.SlowThreshold():
		msg = fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold())
	default:
		msg = ""
	}

	if rows == -1 {
		l.logger.Trace(msg, toto.V{"sql": sql, "elapsed": float64(elapsed.Nanoseconds()) / 1e6, "rows": "-"})
	} else {
		l.logger.Trace(msg, toto.V{"sql": sql, "elapsed": float64(elapsed.Nanoseconds()) / 1e6, "rows": rows})
	}
}
