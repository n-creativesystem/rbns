package contexts

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/internal/contexts"
)

var (
	ErrRedirect = errors.New("redirect")
	ErrIgnore   = errors.New("ignore")
	currentKey  = "contexts.rest"
)

type HandlerFunc = func(c *Context) error

type Context struct {
	*gin.Context

	Tenant     string
	LoginUser  *model.LoginUser
	IsSignedIn bool
	next       bool
}

func (c *Context) Value(key interface{}) interface{} {
	v := c.Context.Value(key)
	if v != nil {
		return v
	}
	return c.Context.Request.Context().Value(key)
}

func (c *Context) Next() {
	c.next = true
}

func (c *Context) Save() {
	c.Context.Set(currentKey, c.copy())
}

func (c *Context) reset(gc *gin.Context, parent *Context) {
	c.Context = gc
	c.Tenant = parent.Tenant
	if parent.LoginUser != nil {
		loginUser := &model.LoginUser{}
		*loginUser = *parent.LoginUser
		c.LoginUser = loginUser
	}
	c.IsSignedIn = parent.IsSignedIn
	c.next = false
}

func (c *Context) copy() *Context {
	cc := &Context{
		IsSignedIn: c.IsSignedIn,
		Tenant:     c.Tenant,
		next:       false,
	}
	if c.LoginUser != nil {
		loginUser := *c.LoginUser
		cc.LoginUser = &loginUser
	}
	return cc
}

func GinWrap(fn HandlerFunc) gin.HandlerFunc {
	return func(gc *gin.Context) {
		tenant := contexts.FromTenantContext(gc)
		var c Context
		if v, ok := gc.Get(currentKey); ok {
			value := v.(*Context)
			c.reset(gc, value)
		} else {
			c = Context{
				Context:    gc,
				Tenant:     tenant,
				LoginUser:  nil,
				IsSignedIn: false,
			}
		}
		err := fn(&c)
		if IsErr(err) {
			var resErrJSON response.ErrorResponse
			if errors.As(err, &resErrJSON) {
				c.Context.AbortWithStatusJSON(resErrJSON.Status, resErrJSON)
			} else {
				c.Context.AbortWithStatusJSON(http.StatusBadRequest, response.ErrJson("Bad request", err))
			}
		} else {
			if c.next {
				c.Save()
				c.Context.Next()
			}
		}
	}
}

func IsErr(err error) bool {
	return (err != nil && !(errors.Is(err, ErrRedirect) || errors.Is(err, ErrIgnore)))
}

func ChangeContext(gc *gin.Context) *Context {
	v, ok := gc.Get(currentKey)
	if !ok {
		return nil
	}
	c, ok := v.(*Context)
	if !ok {
		return nil
	}
	c.Context = gc
	return c
}
