package restserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Copy() Error {
	return *e
}

func NewError(code int, message, description string) *Error {
	return &Error{
		Code:        code,
		Message:     message,
		Description: description,
	}
}

type errorStatus int

const (
	reqParamErr errorStatus = 1001 + iota
	reqBodyErr
	dbErr
	authErr

	systemErr errorStatus = 9999
)

var (
	ErrInvalidRequestParameter = NewError(int(reqParamErr), "Invalid request parameter", "")
	ErrInvalidRequestBody      = NewError(int(reqBodyErr), "Invalid request body", "")
)

type invalidErr int

const (
	body invalidErr = iota
	parameter
	response
)

func errorMiddleware(c *gin.Context) {
	c.Next()
	ginErr := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if ginErr != nil {
		logrus.Error(ginErr.Err)
		switch ginErr.Meta {
		case body:
			errBody := ErrInvalidRequestBody.Copy()
			errBody.Description = ginErr.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, errBody)
		case parameter:
			errBody := ErrInvalidRequestParameter.Copy()
			errBody.Description = ginErr.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, errBody)
		case response:
			var err error
			var errBody *Error
			if st, ok := status.FromError(ginErr.Err); ok {
				switch st.Code() {
				case http.StatusConflict, http.StatusBadRequest, http.StatusNotFound:
					errBody = NewError(int(st.Code()), "data error", st.Message())
					c.AbortWithStatusJSON(int(st.Code()), errBody)
					return
				default:
					err = errors.New(st.Message())
				}
			} else {
				err = ginErr.Err
			}
			if model.IsDefinitionError(err) {
				errBody = NewError(int(dbErr), err.Error(), err.Error())
				c.AbortWithStatusJSON(http.StatusBadRequest, errBody)
			} else {
				switch err.(type) {
				case model.DBErr:
					errBody = NewError(int(dbErr), "database error", err.Error())
				default:
					errBody = NewError(int(systemErr), "system error", err.Error())
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, errBody)
			}
		default:
			errBody := NewError(int(systemErr), "system error", ginErr.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, errBody)
		}
	}
}

func requestError(c *gin.Context, err error, errType invalidErr) bool {
	if err == nil {
		return false
	}
	c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(errType)
	c.Abort()
	return true
}

func responseError(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(response)
		c.Abort()
		return true
	}
	return false
}

func grpcCode2ErrCode(status codes.Code) errorStatus {
	switch status {
	case 2000:
		return authErr
	default:
		return errorStatus(status)
	}
}
