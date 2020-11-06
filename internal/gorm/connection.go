package gorm

import (
	"github.com/totoval/databaser/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connection struct {
	gormDB   *gorm.DB
	configer internal.Configer
}

func (c *Connection) SetConfiger(configer internal.Configer) {
	c.configer = configer
}
func (c *Connection) Connector() internal.Connector {
	return c
}

func (c *Connection) Connect() error {
	var err error
	c.gormDB, err = gorm.Open(mysql.Open(c.configer.Dsn()), c.configer.Build().(*gorm.Config))
	if err != nil {
		return err
	}

	return nil
}
func (c *Connection) Disconnect() error {
	db, err := c.gormDB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
func (c *Connection) Client() interface{} {
	return c.gormDB
}
