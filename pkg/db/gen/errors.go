package db

import (
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// OK is for all OK
	OK = 0

	// UnknownError - errors that is not classified
	UnknownError = 2

	// NoMatchingDocument - error for no documents
	NoMatchingDocument = 47

	// DuplicateKey - Duplicate key
	DuplicateKey = 11000
)

// DBError - Database errors
type DBError struct {
	Msg         string
	Code        int
	WriteErrors []mongo.WriteError
}

func (e *DBError) Error() string {
	if e == nil {
		return "unknown error"
	}
	return e.Msg
}

// ParseError - Parse db error from Error message
func ParseError(err error) *DBError {
	wErr, ok := err.(mongo.WriteException)
	if ok {
		code := UnknownError
		message := ""
		for _, v := range wErr.WriteErrors {
			message = v.Message
			code = v.Code
		}
		return &DBError{
			Msg:         message,
			Code:        code,
			WriteErrors: wErr.WriteErrors,
		}
	} else if i := strings.Index(err.Error(), "mongo: no documents in result"); i >= 0 {
		return &DBError{
			Msg:  "No matching rows",
			Code: NoMatchingDocument,
		}
	}
	return &DBError{
		Msg:  err.Error(),
		Code: UnknownError,
	}
}
