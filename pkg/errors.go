package pkg

import "errors"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrInvalidInput           = errors.New("invalid input")
	ErrAlreadyExists          = errors.New("object exists in the DB")
	ErrWrongCredentials       = errors.New("wrong credentials")
	ErrNoValue                = errors.New("value is missing")
	ErrNonSufficientFunds     = errors.New("non-sufficient funds")
)
