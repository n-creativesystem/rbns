package model

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/n-creativesystem/rbns/utilsconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ErrorStatus struct {
	Code    uint32 `json:"code"`
	Message string `json:"error"`
}

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
	var errMsg ErrorMessage
	if errors.As(err, &errMsg) {
		if errMsg.ID == ErrNoDataFound {
			return true
		}
	}
	return false
}

type i18n struct {
	ja string
	en string
}

type errorBody struct {
	i18n   i18n
	status uint32
}

type errorMessage map[messageID]errorBody

type messageID string

func (m messageID) Error() string {
	return GetErrorMessage(m, en.New()).Error()
}

type ErrorMessage struct {
	ID      messageID
	Status  uint32
	Message string
}

const (
	ErrRequired           messageID = "required"
	ErrNoDataFound        messageID = "no data found"
	ErrAlreadyExists      messageID = "Already exists"
	ErrTenantEmpty        messageID = "Tenant is empty"
	ErrTenantNameValidate messageID = "Tenant name is validation error"
	ErrInternal           messageID = "Internal error"
	ErrForbidden          messageID = "Forbidden"
)

var (
	message = errorMessage{
		ErrForbidden: errorBody{
			i18n: i18n{
				ja: "アクセス権限がありません",
				en: http.StatusText(http.StatusForbidden),
			},
			status: uint32(codes.PermissionDenied),
		},
		ErrRequired: errorBody{
			i18n: i18n{
				ja: "必須フィールが未入力です",
				en: "Required field of empty",
			},
			status: uint32(codes.InvalidArgument),
		},
		ErrNoDataFound: errorBody{
			i18n: i18n{
				ja: "データが見つかりません",
				en: "No data found",
			},
			status: uint32(codes.NotFound),
		},
		ErrAlreadyExists: errorBody{
			i18n: i18n{
				ja: "既に存在しています",
				en: "Already exists",
			},
			status: uint32(codes.AlreadyExists),
		},
		ErrTenantEmpty: errorBody{
			i18n: i18n{
				ja: "テナントが選択されていません",
				en: "No tenant selected",
			},
			status: uint32(codes.InvalidArgument),
		},
		ErrTenantNameValidate: errorBody{
			i18n: i18n{
				ja: "テナント名に使用できない文字や記号が含まれています",
				en: "The tenant name contains characters and symbols that cannot be used",
			},
			status: uint32(codes.InvalidArgument),
		},
		ErrInternal: errorBody{
			status: uint32(codes.Internal),
		},
	}
	messageLock sync.RWMutex
)

func (err ErrorMessage) ToErrorStatus() ErrorStatus {
	return NewErrorStatus(err.Status, err.Message)
}

func (err ErrorMessage) Error() string {
	return err.ToErrorStatus().Error()
}

func GetErrorMessage(id messageID, translator locales.Translator, messages ...string) ErrorMessage {
	messageLock.RLock()
	defer messageLock.RUnlock()
	if v, ok := message[id]; ok {
		message := ""
		switch translator.Locale() {
		case "ja":
			message = v.i18n.ja
		default:
			message = v.i18n.en
		}
		return ErrorMessage{
			ID:      id,
			Status:  v.status,
			Message: message,
		}
	}
	v := message[ErrInternal]
	var buf strings.Builder
	for i, message := range messages {
		if i == 0 {
			buf.WriteString(message)
			continue
		}
		buf.WriteString("\n")
		buf.WriteString(message)
	}
	return ErrorMessage{
		ID:      ErrInternal,
		Status:  v.status,
		Message: buf.String(),
	}
}

func GetErrorStatusWithErr(err error, translator locales.Translator, messages ...string) ErrorStatus {
	var errId messageID
	if errors.As(err, &errId) {
		return GetErrorMessage(errId, translator, messages...).ToErrorStatus()
	}
	return NewErrorStatus(uint32(codes.Internal), err.Error())
}
