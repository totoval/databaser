package internal

import (
	"database/sql"
	"github.com/totoval/databaser/pkg/facade"
	facade_logger "github.com/totoval/logger/pkg/facade"
)

type Databaser interface {
	Connector
	New() Databaser

	Queryer
}

type DryRunner interface {
	// Execute the given callback in "dry run" mode.
	Pretend(callbackFunc func(db facade.Queryer) interface{}) (callbackReturn interface{})
}

type Queryer interface {
	// Table(table string) Builder // Begin a fluent query against a database table.
	SetLogger(logger facade_logger.Logger) error

	SetConnector(connector Connector) // Set DB Driver Client

	Raw(query string, bindings ...interface{}) facade.Scanner                     // Get a new raw query expression. @todo error, all
	SelectOne(query string, bindings ...interface{}) facade.OneScanner            // Run a select statement and return a single result. @todo scan, error , USE READ_PDO
	Select(query string, bindings ...interface{}) facade.Scanner                  // Run a select statement against the database. @todo scan, error , USE READ_PDO
	Insert(query string, bindings ...interface{}) (lastInsertId int64, err error) // Run an insert statement against the database.@todo last insert id, error
	Update(query string, bindings ...interface{}) (affectedRows int64, err error) // Run an update statement against the database. @todo affected rows, error
	Delete(query string, bindings ...interface{}) (affectedRows int64, err error) // Run a delete statement against the database. @todo affected rows, error

	Statement(query string, bindings ...interface{}) (affectedRows int64, err error) // Execute an SQL statement and return the boolean result. @todo  error
	PrepareBindings(bindings ...interface{}) []interface{}                           // Prepare the query bindings for execution. @todo AllDate to string, All bool to int

	Unprepared(query string) facade.Scanner // Run a raw, unprepared query against the PDO connection. @todo error

	Transactioner
	QueryForker
	DryRunner
}

type Transactioner interface {
	Transaction(transactionFunc func(tx facade.Queryer) error, attempts uint) // Execute a Closure within a transaction. Use `panic` inner for error
	TransactionBegin() facade.Queryer                                         // Start a new database transaction.
	TransactionCommit() error                                                 // Commit the active database transaction.
	TransactionRollBack() error                                               // Rollback the active database transaction.
	// TransactionLevel() uint // Get the number of active transactions.
}

type TransactionExceptionHandler interface {
	HandleTransactionException(tx facade.Transactioner, transactionFunc func(tx facade.Queryer) error, err error, currentAttempt uint, maxAttempts uint, log facade_logger.Logger)
}

type QueryForker interface {
	Fork(clientPtr interface{}) facade.Queryer
}

// type Builder interface {
// 	//@todo in eloquent
// }

type Scanner interface {
	OneScanner
	MultiScanner
}

type OneScanner interface {
	Scan(dest interface{}) error
}

type MultiScanner interface {
	ScanRows(rows *sql.Rows, dest interface{}) error
}
