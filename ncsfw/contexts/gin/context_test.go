// package gin_test

// import (
// 	"bytes"
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
// 	"github.com/n-creativesystem/rbns/handler/restserver/middleware/otel"
// 	"github.com/n-creativesystem/rbns/ncsfw/tracer"
// 	"github.com/n-creativesystem/rbns/logger"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRestServer(t *testing.T) {
// 	log := logger.New("contexts test")
// 	ctx := context.Background()
// 	trace_, _ := tracer.InitOpenTelemetry("github.com/n-creativesystem/rbns")
// 	defer trace_.Cleanup(ctx)
// 	var buf bytes.Buffer
// 	r := gin.New()
// 	r.GET("/aaa",
// 		contexts.GinWrap(otel.Middleware("rest server")),
// 		otel.RestLogger(log),
// 		contexts.GinWrap(func(c *contexts.Context) error {
// 			ctx, span := tracer.Start(c, "tenant set")
// 			defer span.End()
// 			buf.WriteString("a")
// 			c.Tenant = "tenant"
// 			c.Request = c.Request.WithContext(ctx)
// 			c.Save()
// 			c.Context.Next()
// 			buf.WriteString("b")
// 			return nil
// 		}),
// 		func(c *gin.Context) {
// 			buf.WriteString("c")
// 			ctx := contexts.ChangeContext(c)
// 			assert.Equal(t, "tenant", ctx.Tenant)
// 			c.Next()
// 		},
// 		contexts.GinWrap(func(c *contexts.Context) error {
// 			buf.WriteString("d")
// 			c.Next()
// 			return nil
// 		}),
// 		contexts.GinWrap(func(c *contexts.Context) error {
// 			buf.WriteString("e")
// 			c.JSON(http.StatusOK, gin.H{})
// 			return nil
// 		}))
// 	var w *httptest.ResponseRecorder
// 	var req *http.Request
// 	w = httptest.NewRecorder()
// 	req = httptest.NewRequest(http.MethodGet, "/aaa", bytes.NewBufferString(""))
// 	req.Header.Add("x-context-test", "test")

// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, "acdeb", buf.String())
// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// }
