package model

import (
	"net/http"

	"google.golang.org/grpc/status"
)

type ModelError struct {
	err error
}

func (e ModelError) Error() string {
	return e.err.Error()
}

func IsDefinitionError(err error) bool {
	switch err {
	case ErrRequired, ErrNoData:
		return true
	default:
		switch err.(type) {
		case ModelError:
			return true
		default:
			return false
		}
	}
}

var (
	// ErrRequired `Required field of empty`
	ErrRequired = status.Error(http.StatusBadRequest, "Required field of empty")
	// ErrNoData `No data found`
	ErrNoData = status.Error(http.StatusNotFound, "No data found")
	// ErrAlreadyExists `Already exists`
	ErrAlreadyExists = status.Error(http.StatusConflict, "Already exists")
)

type DBErr struct {
	err error
}

func (e DBErr) Error() string {
	return e.err.Error()
}

func NewDBErr(err error) error {
	if err == nil {
		return nil
	}
	return DBErr{
		err: err,
	}
}
