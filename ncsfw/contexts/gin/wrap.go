// package gin

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	. "github.com/n-creativesystem/rbns/handler/restserver/contexts"
// 	"github.com/n-creativesystem/rbns/handler/restserver/response"
// 	"github.com/n-creativesystem/rbns/ncsfw/logger"
// )

// func GinWrap(fn HandlerFunc) gin.HandlerFunc {
// 	return func(gc *gin.Context) {
// 		tenant := FromTenantContext(gc.Request.Context())
// 		var c context
// 		if v, ok := gc.Get(currentKey); ok {
// 			value := v.(*context)
// 			c.reset(gc, value)
// 		} else {
// 			c = context{
// 				gc:         gc,
// 				Tenant:     tenant,
// 				LoginUser:  nil,
// 				IsSignedIn: false,
// 				log:        logger.New("contexts.gin"),
// 			}
// 		}
// 		err := fn(&c)
// 		if IsErr(err) {
// 			var resErrJSON response.ErrorResponse
// 			if errors.As(err, &resErrJSON) {
// 				c.gc.AbortWithStatusJSON(resErrJSON.Status, resErrJSON)
// 			} else {
// 				c.gc.AbortWithStatusJSON(http.StatusBadRequest, response.ErrJson("Bad request", err))
// 			}
// 			c.log.ErrorWithContext(c.Request().Context(), err, "request error")
// 		} else {
// 			if c.next {
// 				c.Save()
// 				c.gc.Next()
// 			}
// 		}
// 	}
// }

// func IsErr(err error) bool {
// 	return (err != nil && !(errors.Is(err, ErrRedirect) || errors.Is(err, ErrIgnore)))
// }

// func ChangeContext(gc *gin.Context) Context {
// 	v, ok := gc.Get(currentKey)
// 	if !ok {
// 		return nil
// 	}
// 	c, ok := v.(*context)
// 	if !ok {
// 		return nil
// 	}
// 	c.gc = gc
// 	return c
// }
