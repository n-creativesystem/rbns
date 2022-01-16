package restserver

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/middleware/auth"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/ncsfw"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/middleware"
	"github.com/n-creativesystem/rbns/ncsfw/middleware/otel"
	"github.com/n-creativesystem/rbns/ncsfw/tenants"
	"github.com/n-creativesystem/rbns/service"
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
	authMiddleware *auth.AuthMiddleware
	store          sessions.Store
	tenantService  service.Tenant
	apiKeyService  service.APIKey

	Cfg *config.Config
}

func (hs *HTTPServer) registerRouting() {
	authService := hs.authMiddleware.RestMiddleware(hs.store)
	r := ncsfw.New()
	r.Delims("[[", "]]")
	r.Use(otel.Middleware("rest server"), middleware.Recover())
	r.LoadHTMLGlob(path.Join(hs.Cfg.StaticFilePath, "/*.html"))
	r.Static("/static", hs.Cfg.StaticFilePath)
	login := r.Group("login")
	{
		login.Get("/:name", hs.OAuthLogin, otel.Middleware("login request"))
		login.Get("/provider", hs.GetOAuthProvider, otel.Middleware("login provider request"))
	}
	r.Get("/logout", hs.Logout, otel.Middleware("logout request"))

	g := r.Group("/api", authService.Required, otel.Middleware("api request"))
	{
		requiredTenant := g.Group("", tenants.HTTPServerMiddleware())
		{
			requiredTenant.Any("/v1/g/*gateway", hs.WrapGateway(hs.gateway), otel.Middleware("grpc gateway request"))
			requiredTenant.Post("/auth/keys", hs.AddAPIKey, otel.Middleware("generate api key request"), authService.RoleCheck(model.ROLE_ADMIN))
			requiredTenant.Delete("/auth/keys/:id", hs.DeleteAPIKey, otel.Middleware("delete api keys request"), authService.RoleCheck(model.ROLE_ADMIN))
		}
		tenant := g.Group("/v1/tenants", otel.Middleware("tenant request"))
		{
			tenant.Get("", hs.getTenants, otel.Middleware("add tenant request"), authService.RoleCheck(model.ROLE_VIEWER))
			tenant.Post("", hs.addTenant, otel.Middleware("get tenants request"), authService.RoleCheck(model.ROLE_ADMIN))
		}
	}
	r.NoRoute(hs.Index, authService.NotRequired)
	hs.handler = r
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

	if !strings.HasPrefix(to.Path, "/") {
		return ErrForbiddenRedirectTo
	}
	if strings.HasPrefix(to.Path, "//") {
		return ErrForbiddenRedirectTo
	}

	rootURL := hs.Cfg.RootURL.String()
	if rootURL != "" && !strings.HasPrefix(to.Path, rootURL+"/") {
		return ErrInvalidRedirectTo
	}

	return nil
}
