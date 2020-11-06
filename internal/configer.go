package internal

import (
	facade_logger "github.com/totoval/logger/pkg/facade"
	"time"
)

type Configer interface {
	SetLogger(logger facade_logger.Logger) error
	Logger() facade_logger.Logger
	SetSlowSqlThreshold(threshold time.Duration) error
	SlowSqlThreshold() (threshold time.Duration)
	SetNowFunc(nowFunc func() time.Time) error
	NowFunc() func() time.Time
	SetDryRun(dryRun bool) error
	SetDsn(dsn string) error
	Dsn() string
	Build() interface{}
	Configer() Configer
}
