package restserver

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/internal/saml"
	"github.com/n-creativesystem/rbns/logger"
)

const (
	samlKey = "x-saml-payload"
)

func samlHandler(r *gin.Engine, opts ...saml.Options) {
	middleware := saml.New(opts...)
	handler := func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		middleware.ServeHTTP(w, r)
	}
	r.Any("saml/*actions", handler)
	r.Use(func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		log := logger.FromContext(r.Context())
		session, err := middleware.Session.GetSession(r)
		if err != nil {
			log.Errorf("Session error: %s", err)
		}
		if session != nil {
			buf, err := json.Marshal(session)
			if err != nil {
				log.Errorf("Session Json Marshal: %s", err)
			}
			strJwt := base64.RawURLEncoding.EncodeToString(buf)
			c.Request = r.WithContext(samlsp.ContextWithSession(r.Context(), session))
			w.Header().Add(samlKey, strJwt)
			c.Set(samlKey, strJwt)
			c.Next()
			return
		}
		if err == samlsp.ErrNoSession {
			status := w.Status()
			middleware.HandleStartAuthFlow(w, r)
			if w.Status() != status {
				return
			}
			c.Next()
			return
		}
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	})
}
