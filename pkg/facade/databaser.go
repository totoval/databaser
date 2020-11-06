package facade

import (
	"database/sql"
	facade_logger "github.com/totoval/logger/pkg/facade"
)

type Databaser interface {
	Queryer
}

type Queryer interface {
	Raw(query string, bindings ...interface{}) Scanner                            // Get a new raw query expression. @todo error, all
	SelectOne(query string, bindings ...interface{}) OneScanner                   // Run a select statement and return a single result. @todo scan, error , USE READ_PDO
	Select(query string, bindings ...interface{}) Scanner                         // Run a select statement against the database. @todo scan, error , USE READ_PDO
	Insert(query string, bindings ...interface{}) (lastInsertId int64, err error) // Run an insert statement against the database.@todo last insert id, error
	Update(query string, bindings ...interface{}) (affectedRows int64, err error) // Run an update statement against the database. @todo affected rows, error
	Delete(query string, bindings ...interface{}) (affectedRows int64, err error) // Run a delete statement against the database. @todo affected rows, error

	Statement(query string, bindings ...interface{}) (affectedRows int64, err error) // Execute an SQL statement and return the boolean result. @todo  error
	PrepareBindings(bindings ...interface{}) []interface{}                           // Prepare the query bindings for execution. @todo AllDate to string, All bool to int

	Unprepared(query string) Scanner // Run a raw, unprepared query against the PDO connection. @todo error

	Transactioner
	QueryForker
	DryRunner
}

type DryRunner interface {
	Pretend(callbackFunc func(db Queryer) interface{}) (callbackReturn interface{})
}

type QueryForker interface {
	Fork(clientPtr interface{}) Queryer
}

type Transactioner interface {
	Transaction(transactionFunc func(tx Queryer) error, attempts uint) // Execute a Closure within a transaction. Use `panic` inner for error
	TransactionBegin() Queryer                                         // Start a new database transaction.
	TransactionCommit() error                                          // Commit the active database transaction.
	TransactionRollBack() error                                        // Rollback the active database transaction.

	HandleTransactionException(tx Transactioner, transactionFunc func(tx Queryer) error, err error, currentAttempt uint, maxAttempts uint, log facade_logger.Logger)
}

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
