package errors

import "fmt"

type InternalError struct {
	cause string
}

func NewInternalError(cause error) *InternalError {
	return &InternalError{
		cause: cause.Error(),
	}
}

func (e *InternalError) Error() string {
	return e.cause
}

type DataError struct {
	message string
}

func NewDataError(message string) *DataError {
	return &DataError{
		message: message,
	}
}

func (e *DataError) Error() string {
	return e.message
}

type AccountNotFoundError struct {
	accountID string
}

func NewAccountNotFoundError(accountID string) *AccountNotFoundError {
	return &AccountNotFoundError{
		accountID: accountID,
	}
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("account %s was not found in system", e.accountID)
}
