package model

import (
	"errors"

	"github.com/n-creativesystem/rbns/utilsconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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

type ErrorStatus struct {
	Code    uint32 `json:"code"`
	Message string `json:"error"`
}

var (
	// ErrRequired `Required field of empty`
	ErrRequired = NewErrorStatus(uint32(codes.InvalidArgument), "Required field of empty")
	// ErrNoData `No data found`
	ErrNoData = NewErrorStatus(uint32(codes.NotFound), "No data found")
	// ErrAlreadyExists `Already exists`
	ErrAlreadyExists = NewErrorStatus(uint32(codes.AlreadyExists), "Already exists")

	ErrTenantRequired = NewErrorStatus(uint32(codes.InvalidArgument), "Tenant is empty")
)

func NewErrorStatus(code uint32, message string) ErrorStatus {
	return ErrorStatus{
		Code:    code,
		Message: message,
	}
}

func (e ErrorStatus) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Code), e.Message)
}

func (e ErrorStatus) REST() error {
	return ErrorStatus{
		Code:    uint32(utilsconv.GRPCStatus2HTTPStatus(codes.Code(e.Code))),
		Message: e.Message,
	}
}

func (e ErrorStatus) Error() string {
	return e.GRPCStatus().String()
}

func IsNoData(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if errors.Is(err, ErrNoData) {
		return true
	}
	return false
}
