package contexts

// import (
// 	ctx "context"
// 	"errors"
// 	"net/http"

// 	"github.com/n-creativesystem/rbns/domain/model"
// 	"github.com/n-creativesystem/rbns/ncsfw/logger"
// )

// var (
// 	ErrRedirect = errors.New("redirect")
// 	ErrIgnore   = errors.New("ignore")
// 	currentKey  = "contexts.rest"
// )

// type context struct {
// 	writer     *responseWriter
// 	request    *http.Request
// 	Tenant     string
// 	LoginUser  *model.LoginUser
// 	IsSignedIn bool
// 	next       bool
// 	log        logger.Logger
// 	param      Params
// }

// var _ Context = (*context)(nil)

// func (c *context) BaseContext() interface{} {
// 	return nil
// }

// func (c *context) Value(key interface{}) interface{} {
// 	return c.request.Context().Value(key)
// }

// func (c *context) Next() {
// 	c.next = true
// }

// func (c *context) Save() {
// 	*c.request = *c.request.WithContext(ctx.WithValue(c.request.Context(), currentKey, c.copy()))
// }

// func (c *context) copy() *context {
// 	cc := &context{
// 		writer:     &responseWriter{},
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
// 	return c.request
// }

// func (c *context) JSON(code int, obj interface{}) {
// 	c.JSON(code, obj)
// }

// func (c *context) Writer() ResponseWriter {
// 	return c.writer
// }

// func (c *context) SetRequest(r *http.Request) {
// 	*c.request = *r
// }

// func (c *context) Errors() []error {
// 	return nil
// 	// err := make([]error, 0, 10)
// 	// for _, e := range c.gc.Errors {
// 	// 	err = append(err, e)
// 	// }
// 	// return err
// }

// func (c *context) FullPath() string {
// 	return ""
// }

// func (c *context) IsAborted() bool {
// 	return false
// }
// func (c *context) Abort() {
// 	// c.gc.Abort()
// }
// func (c *context) AbortWithStatus(code int) {
// 	// c.gc.AbortWithStatus(code)
// }
// func (c *context) AbortWithStatusJSON(code int, jsonObj interface{}) {
// 	// c.gc.AbortWithStatusJSON(code, jsonObj)
// }
// func (c *context) AbortWithError(code int, err error) error {
// 	// return c.gc.AbortWithError(code, err)
// 	return nil
// }

// func (c *context) Logger() logger.Logger {
// 	return c.log
// }

// func (c *context) File(filePath string) {
// 	http.ServeFile(c.writer, c.request, filePath)
// }

// func (c *context) Param(param string) string {
// 	return c.param.ByName(param)
// }
