package gorm

import (
	"github.com/totoval/databaser/internal"
)

type Database struct {
	Connection
	Query
}

func (d *Database) New() internal.Databaser {
	d.Connection = Connection{}
	return d
}
