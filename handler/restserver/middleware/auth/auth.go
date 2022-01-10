package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/config"
)

type Middleware interface {
	Use(r *gin.Engine)
}

const AuthKey string = "authKey"

type noAuth struct{}

func (n *noAuth) middleware(c *gin.Context) {
	c.Set(AuthKey, "")
	c.Next()
}

func (n *noAuth) Use(r *gin.Engine) {
	r.Use(n.middleware)
}

func New(conf *config.Config, store sessions.Store) (Middleware, error) {
	// if conf.IsSAML() {
	// 	return NewSAML(conf, store)
	// }
	// if conf.IsOIDC() {
	// 	return NewOIDC(conf, store)
	// }

	return &noAuth{}, nil
}
