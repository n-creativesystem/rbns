package ncsfw

import (
	"net/http"

	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type HandlerFunc func(c Context) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Context interface {
	Reset(w http.ResponseWriter, r *http.Request, fullPath string, params *Params)
	Request() *http.Request
	SetRequest(*http.Request)
	Writer() ResponseWriter
	File(filePath string)
	Param(param string) string

	SetTenant(tenant string)
	GetTenant() string
	SetLoginUser(loginUser interface{})
	GetLoginUser() interface{}
	IsSignedIn() bool
	SignedIn(isSigned bool)

	Logger() logger.Logger
	FullPath() string
	ClientIP() string

	Status(code int)
	HTML(code int, name string, obj interface{})
	JSON(code int, obj interface{})
	Redirect(code int, url string)

	BindJSON(obj interface{}) error
	Query(key string) string
	GetHeader(key string) string
}
