package gorm

import (
	"github.com/totoval/databaser/internal"
	facade_logger "github.com/totoval/logger/pkg/facade"
	"gorm.io/gorm"
	"time"
)

type Configuration struct {
	logger        facade_logger.Logger
	slowThreshold time.Duration
	dsn           string
	nowFunc       func() time.Time
	dryRun        bool
}

func (c *Configuration) Configer() internal.Configer {
	return c
}
func (c *Configuration) Build() interface{} {
	return &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 NewLog(c.logger, c.SlowSqlThreshold()),
		NowFunc:                c.NowFunc(),
		DryRun:                 c.dryRun,
	}
}
func (c *Configuration) SetDsn(dsn string) error {
	c.dsn = dsn
	return nil
}
func (c *Configuration) Dsn() string {
	return c.dsn
}
func (c *Configuration) SetLogger(logger facade_logger.Logger) error {
	c.logger = logger
	return nil
}
func (c *Configuration) Logger() facade_logger.Logger {
	return c.logger
}

func (c *Configuration) SetSlowSqlThreshold(threshold time.Duration) error {
	c.slowThreshold = threshold
	return nil
}
func (c *Configuration) SlowSqlThreshold() (threshold time.Duration) {
	return c.slowThreshold
}

func (c *Configuration) SetNowFunc(nowFunc func() time.Time) error {
	c.nowFunc = nowFunc
	return nil
}

func (c *Configuration) NowFunc() func() time.Time {
	if c.nowFunc == nil {
		return time.Now
	}
	return c.nowFunc
}

func (c *Configuration) SetDryRun(dryRun bool) error {
	c.dryRun = dryRun
	return nil
}
