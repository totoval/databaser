package gorm

import (
	"errors"
	"fmt"
	"github.com/totoval/databaser/internal"
	"github.com/totoval/databaser/internal/common"
	"github.com/totoval/databaser/pkg/facade"
	facade_logger "github.com/totoval/logger/pkg/facade"
	"gorm.io/gorm"
)

type Query struct {
	common.QueryPrepareBindings
	common.TransactionExceptionHandler
	connector internal.Connector
	client    *gorm.DB
	logger    facade_logger.Logger
}

func (q *Query) SetLogger(logger facade_logger.Logger) error {
	q.logger = logger
	return nil
}

func (q *Query) SetConnector(connector internal.Connector) {
	q.connector = connector
	q.client = q.connector.Client().(*gorm.DB)
}
func (q *Query) Raw(query string, bindings ...interface{}) facade.Scanner {
	return &Scan{gormDB: q.client.Raw(query, bindings...)}
}

func (q *Query) SelectOne(query string, bindings ...interface{}) facade.OneScanner {
	return q.Raw(query, bindings...)
}

func (q *Query) Select(query string, bindings ...interface{}) facade.Scanner {
	return q.Raw(query, bindings...)
}

func (q *Query) Insert(query string, bindings ...interface{}) (lastInsertId int64, err error) {
	db, err := q.client.DB()
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(query, bindings...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (q *Query) Update(query string, bindings ...interface{}) (affectedRows int64, err error) {
	return q.Statement(query, bindings...)
}

func (q *Query) Delete(query string, bindings ...interface{}) (affectedRows int64, err error) {
	return q.Statement(query, bindings...)
}

func (q *Query) Statement(query string, bindings ...interface{}) (affectedRows int64, err error) {
	result := q.client.Exec(query, q.PrepareBindings(bindings...)...)
	err = result.Error
	if err != nil {
		return 0, err
	}
	if result.RowsAffected < 0 {
		affectedRows = 0
	}
	return affectedRows, nil
}

func (q *Query) Unprepared(query string) facade.Scanner {
	return q.Raw(query)
}

func (q *Query) Transaction(transactionFunc func(tx facade.Queryer) error, attempts uint) {
	if attempts <= 0 {
		attempts = 1
	}
	currentAttempt := uint(1)

	tx := q.TransactionBegin()
	defer func(_tx facade.Queryer) {
		if err := recover(); err != nil {
			var __err error
			if _err, ok := err.(error); ok {
				__err = _err
			} else {
				__err = errors.New(fmt.Sprint(err))
			}
			tx.(internal.TransactionExceptionHandler).HandleTransactionException(tx, transactionFunc, __err, currentAttempt, attempts, q.logger)
		}
	}(tx)

	if err := transactionFunc(tx); err != nil {
		q.logger.Panic(err)
	}

	if err := q.TransactionCommit(); err != nil {
		q.logger.Panic(err)
	}
}
func (q *Query) Fork(clientPtr interface{}) facade.Queryer {
	return &Query{
		QueryPrepareBindings:        q.QueryPrepareBindings,
		TransactionExceptionHandler: q.TransactionExceptionHandler,
		connector:                   nil, // set to nil, for safe
		client:                      clientPtr.(*gorm.DB),
		logger:                      q.logger,
	}
}
func (q *Query) TransactionBegin() facade.Queryer {
	return q.Fork(q.client.Begin())
}

func (q *Query) TransactionCommit() error {
	return q.client.Commit().Error
}

func (q *Query) TransactionRollBack() error {
	return q.client.Rollback().Error
}

func (q *Query) Pretend(callbackFunc func(db facade.Queryer) interface{}) (callbackReturn interface{}) {
	return callbackFunc(q.Fork(q.client.Session(&gorm.Session{DryRun: true})))
}
