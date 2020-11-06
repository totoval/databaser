package main

import (
	"github.com/totoval/databaser/pkg"
	"github.com/totoval/databaser/pkg/facade"
	"github.com/totoval/framework/helpers/toto"
	pkg_logger "github.com/totoval/logger/pkg"
	facade_logger "github.com/totoval/logger/pkg/facade"
	structs_logger "github.com/totoval/logger/pkg/structs"
	"time"
)

//@todo register into global vars
var log facade_logger.Logger
var db facade.Databaser

func main() {
	// use driver then config
	l := &pkg_logger.Log{}
	if err := l.Use(pkg_logger.DriverLogrus).Config(map[string]interface{}{
		"level": structs_logger.LevelTrace,
	}); err != nil {
		panic(err)
	}
	// get the facade
	log = l.Component().(facade_logger.Logger)

	log.Info("test", toto.V{"haha": 123}, toto.V{"1": 2})

	d := &pkg.Database{}
	if err := d.Use(pkg.DriverGorm).Config(map[string]interface{}{
		"logger":                log,
		"dsn":                   "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		"slow_sql_threshold_ms": time.Duration(5000),
		"dry_run":               false,
		"now_func":              time.Now,
	}); err != nil {
		panic(err)
	}
	// get the facade
	db = d.Component().(facade.Databaser)

	data := make(map[string]interface{})
	err := db.Raw("select * from test_table limit 1").Scan(&data)
	if err != nil {
		_ = log.Error(err)
	}
	log.Info(data)
}
