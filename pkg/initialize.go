package pkg

import (
	"errors"
	"github.com/totoval/databaser/internal"
	"github.com/totoval/databaser/internal/gorm"
	facade_logger "github.com/totoval/logger/pkg/facade"
	"time"
)

type Database struct {
	configer  internal.Configer
	connector internal.Connector
	database  internal.Databaser
}

func (d *Database) Component() interface{} {
	return d.database
}
func (d *Database) Use(driver string) Componentor {
	switch driver {
	case DriverGorm:
		d.configer = &gorm.Configuration{}
		d.connector = &gorm.Connection{}
		d.database = &gorm.Database{}
	default:
		d.configer = &gorm.Configuration{}
		d.connector = &gorm.Connection{}
		d.database = &gorm.Database{}
	}
	return d
}

func (d *Database) Config(configuration map[string]interface{}) error {
	d.database.New()

	logger, ok := configuration["logger"].(facade_logger.Logger)
	if !ok {
		return errors.New("unknown configuration logger")
	}
	if err := d.configer.SetLogger(logger); err != nil {
		return err
	}

	slowSqlThresholdMs, ok := configuration["slow_sql_threshold_ms"].(time.Duration)
	if !ok {
		return errors.New("unknown configuration slow_sql_threshold_ms")
	}
	if err := d.configer.SetSlowSqlThreshold(slowSqlThresholdMs * time.Millisecond); err != nil {
		return err
	}

	dryRun, ok := configuration["dry_run"].(bool)
	if !ok {
		return errors.New("unknown configuration dry_run")
	}
	if err := d.configer.SetDryRun(dryRun); err != nil {
		return err
	}

	nowFunc, ok := configuration["now_func"].(func() time.Time)
	if !ok {
		return errors.New("unknown configuration now_func")
	}
	if err := d.configer.SetNowFunc(nowFunc); err != nil {
		return err
	}

	dsn, ok := configuration["dsn"].(string)
	if !ok {
		return errors.New("unknown configuration dsn")
	}
	if err := d.configer.SetDsn(dsn); err != nil {
		return err
	}

	d.connector.SetConfiger(d.configer)
	if err := d.connector.Connect(); err != nil {
		return err
	}
	d.database.SetConnector(d.connector)

	return nil
}
