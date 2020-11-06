package gorm

import (
	"database/sql"
	"gorm.io/gorm"
)

type Scan struct {
	gormDB *gorm.DB
}

func (s *Scan) Scan(dest interface{}) error {
	return s.gormDB.Scan(dest).Error
}

func (s *Scan) ScanRows(rows *sql.Rows, dest interface{}) error {
	return s.gormDB.ScanRows(rows, dest)
}
