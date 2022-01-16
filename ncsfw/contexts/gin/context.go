// package gin

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/n-creativesystem/rbns/domain/model"
// 	. "github.com/n-creativesystem/rbns/handler/restserver/contexts"
// 	"github.com/n-creativesystem/rbns/ncsfw/logger"
// )

// var (
// 	ErrRedirect = errors.New("redirect")
// 	ErrIgnore   = errors.New("ignore")
// 	currentKey  = "contexts.rest"
// )

// type context struct {
// 	gc *gin.Context

// 	Tenant     string
// 	LoginUser  *model.LoginUser
// 	IsSignedIn bool
// 	next       bool
// 	log        logger.Logger
// }

// var _ Context = (*context)(nil)

// func (c *context) BaseContext() interface{} {
// 	return c.gc
// }

// func (c *context) Value(key interface{}) interface{} {
// 	v := c.gc.Value(key)
// 	if v != nil {
// 		return v
// 	}
// 	return c.gc.Request.Context().Value(key)
// }

// func (c *context) Next() {
// 	c.next = true
// }

// func (c *context) Save() {
// 	c.gc.Set(currentKey, c.copy())
// }

// func (c *context) reset(gc *gin.Context, parent *context) {
// 	c.gc = gc
// 	c.Tenant = parent.Tenant
// 	if parent.LoginUser != nil {
// 		loginUser := &model.LoginUser{}
// 		*loginUser = *parent.LoginUser
// 		c.LoginUser = loginUser
// 	}
// 	c.IsSignedIn = parent.IsSignedIn
// 	c.next = false
// }

// func (c *context) copy() *context {
// 	cc := &context{
// 		IsSignedIn: c.IsSignedIn,
// 		Tenant:     c.Tenant,
// 		next:       false,
// 	}
// 	if c.LoginUser != nil {
// 		loginUser := *c.LoginUser
// 		cc.LoginUser = &loginUser
// 	}
// 	return cc
// }

// func (c *context) Request() *http.Request {
// 	return c.gc.Request
// }

// func (c *context) JSON(code int, obj interface{}) {
// 	c.gc.JSON(code, obj)
// }

// func (c *context) Writer() ResponseWriter {
// 	return c.gc.Writer
// }

// func (c *context) SetRequest(r *http.Request) {
// 	*c.gc.Request = *r
// }

// func (c *context) Errors() []error {
// 	err := make([]error, 0, 10)
// 	for _, e := range c.gc.Errors {
// 		err = append(err, e)
// 	}
// 	return err
// }

// func (c *context) FullPath() string {
// 	return c.gc.FullPath()
// }

// func (c *context) IsAborted() bool {
// 	return c.gc.IsAborted()
// }
// func (c *context) Abort() {
// 	c.gc.Abort()
// }
// func (c *context) AbortWithStatus(code int) {
// 	c.gc.AbortWithStatus(code)
// }
// func (c *context) AbortWithStatusJSON(code int, jsonObj interface{}) {
// 	c.gc.AbortWithStatusJSON(code, jsonObj)
// }
// func (c *context) AbortWithError(code int, err error) error {
// 	return c.gc.AbortWithError(code, err)
// }

// func (c *context) Logger() logger.Logger {
// 	return c.log
// }
