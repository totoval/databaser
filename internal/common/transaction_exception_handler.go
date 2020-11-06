package common

import (
	"github.com/totoval/databaser/pkg/facade"
	facade_logger "github.com/totoval/logger/pkg/facade"
)

type TransactionExceptionHandler string

func (h TransactionExceptionHandler) HandleTransactionException(tx facade.Transactioner, transactionFunc func(tx facade.Queryer) error, err error, currentAttempt uint, maxAttempts uint, logger facade_logger.Logger) {
	if err := tx.TransactionRollBack(); err != nil {
		logger.Panic(err)
	}
	if currentAttempt < maxAttempts {
		tx.Transaction(transactionFunc, maxAttempts-currentAttempt)
	}

	logger.Panic(err)
}
