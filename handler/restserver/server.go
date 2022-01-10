package restserver

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/domain/dtos"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/handler/restserver/middleware"
	"github.com/n-creativesystem/rbns/handler/restserver/middleware/otel"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/logger"
)

var (
	ErrEmailNotAllowed       = errors.New("required email domain not fulfilled")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrNoEmail               = errors.New("login provider didn't return an email address")
	ErrProviderDeniedRequest = errors.New("login provider denied login request")
	ErrTooManyLoginAttempts  = errors.New("too many consecutive incorrect login attempts for user - login for user temporarily blocked")
	ErrPasswordEmpty         = errors.New("no password provided")
	ErrUserDisabled          = errors.New("user is disabled")
	ErrAbsoluteRedirectTo    = errors.New("absolute URLs are not allowed for redirect_to cookie value")
	ErrInvalidRedirectTo     = errors.New("invalid redirect_to cookie value")
	ErrForbiddenRedirectTo   = errors.New("forbidden redirect_to cookie value")
)

type HTTPServer struct {
	log            logger.Logger
	handler        http.Handler
	gateway        http.Handler
	socialService  social.Service
	authMiddleware *middleware.AuthMiddleware
	store          sessions.Store

	Cfg *config.Config
}

func (hs *HTTPServer) registerRouting() {
	authService := hs.authMiddleware.RestMiddleware(hs.store)
	gin.DefaultWriter = hs.log
	r := gin.New()
	r.Use(contexts.GinWrap(otel.Middleware("rest server")), otel.RestLogger(hs.log), gin.Recovery(), contexts.GinWrap(otel.Middleware("rest server")))
	r.LoadHTMLGlob(path.Join(hs.Cfg.StaticFilePath, "/*.html"))
	r.Static("/static", hs.Cfg.StaticFilePath)
	login := r.Group("login")
	{
		login.GET("/:name", contexts.GinWrap(otel.Middleware("login request")), contexts.GinWrap(hs.OAuthLogin))
		login.GET("/provider", contexts.GinWrap(otel.Middleware("login provider request")), contexts.GinWrap(hs.GetOAuthProvider))
	}
	r.GET("/logout", contexts.GinWrap(otel.Middleware("logout request")), contexts.GinWrap(hs.Logout))

	g := r.Group("/api", contexts.GinWrap(otel.Middleware("api request")), contexts.GinWrap(authService.Required))
	{
		g.Any("/g/*gateway", contexts.GinWrap(otel.Middleware("grpc gateway request")), gin.WrapH(hs.gateway))
		g.POST("/auth/keys", contexts.GinWrap(otel.Middleware("generate api key request")), contexts.GinWrap(hs.AddAPIKey))
		g.DELETE("/auth/keys/:id", contexts.GinWrap(otel.Middleware("delete api keys request")), contexts.GinWrap(hs.DeleteAPIKey))
		tenant := g.Group("/tenant", contexts.GinWrap(otel.Middleware("tenant request")))
		{
			tenant.GET("", contexts.GinWrap(otel.Middleware("add tenant request")), contexts.GinWrap(hs.getTenants))
			tenant.POST("", contexts.GinWrap(otel.Middleware("get tenants request")), contexts.GinWrap(hs.addTenant))
		}
	}
	r.NoRoute(contexts.GinWrap(authService.NotRequired), contexts.GinWrap(hs.NoRoute))
	hs.handler = r
}

func (s *HTTPServer) NoRoute(c *contexts.Context) error {
	loginUser := &model.LoginUser{}
	if c.IsSignedIn {
		*loginUser = *c.LoginUser
	}
	if !loginUser.Valid() || !loginUser.IsVerify() {
		loginUser = &model.LoginUser{}
		c.IsSignedIn = false
	}
	currentUser := dtos.CurrentUser{
		ID:         loginUser.ID,
		UseName:    loginUser.UseName,
		Email:      loginUser.Email,
		Role:       loginUser.Role,
		Groups:     loginUser.Groups,
		IsSignedIn: c.IsSignedIn,
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"BaseURL": s.Cfg.RootURL.String(),
		"SubPath": s.Cfg.SubPath,
		"User":    currentUser,
	})
	return nil
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

type Redirect interface {
	Redirect(code int, url string)
}

func (s *HTTPServer) Redirect(r Redirect, url string) {
	r.Redirect(http.StatusFound, url)
}

func (hs *HTTPServer) ValidateRedirectTo(redirectTo string) error {
	to, err := url.Parse(redirectTo)
	if err != nil {
		return ErrInvalidRedirectTo
	}
	if to.IsAbs() {
		return ErrAbsoluteRedirectTo
	}

	if to.Host != "" {
		return ErrForbiddenRedirectTo
	}

	// path should have exactly one leading slash
	if !strings.HasPrefix(to.Path, "/") {
		return ErrForbiddenRedirectTo
	}
	if strings.HasPrefix(to.Path, "//") {
		return ErrForbiddenRedirectTo
	}

	// when using a subUrl, the redirect_to should start with the subUrl (which contains the leading slash), otherwise the redirect
	// will send the user to the wrong location
	rootURL := hs.Cfg.RootURL.String()
	if rootURL != "" && !strings.HasPrefix(to.Path, rootURL+"/") {
		return ErrInvalidRedirectTo
	}

	return nil
}
