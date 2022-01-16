package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/metadata"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/internal/utils"
	"github.com/n-creativesystem/rbns/ncsfw"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"github.com/n-creativesystem/rbns/service"
)

const (
	loginKey         = "loginUser"
	OauthSessionName = "oauth_session"
)

var (
	InValidAPIKey = errors.New("invalid API key")
)

type AuthMiddleware struct {
	Cfg           *config.Config
	SocialService social.Service
	ApiService    service.APIKey
	log           logger.Logger
}

type AuthService struct {
	Cfg           *config.Config
	SocialService social.Service
	ApiService    service.APIKey
	session       sessions.Store
	log           logger.Logger
}

func (a *AuthService) middleware(c ncsfw.Context, fn func(c ncsfw.Context) error) error {
	r := c.Request()
	ctx, span := tracer.Start(r.Context(), "Auth - Middleware")
	defer span.End()
	c.SetRequest(r.WithContext(ctx))
	return fn(c)
}

func (a *AuthService) Required(next ncsfw.HandlerFunc) ncsfw.HandlerFunc {
	return func(c ncsfw.Context) error {
		return a.middleware(c, func(c ncsfw.Context) error {
			returnErr := func(err error) bool {
				return err != nil && !(errors.Is(err, ncsfw.ErrIgnore) || errors.Is(err, ncsfw.ErrRedirect))
			}

			if err := a.initContextWithApiKey(c); returnErr(err) {
				return err
			}

			if err := a.initContextWithToken(c); returnErr(err) {
				return err
			}
			if c.IsSignedIn() {
				return next(c)
			} else {
				return response.ErrJsonWithStatus(http.StatusUnauthorized, "no login", nil)
			}
		})
	}
}

func (a *AuthService) NotRequired(next ncsfw.HandlerFunc) ncsfw.HandlerFunc {
	return func(c ncsfw.Context) error {
		return a.middleware(c, func(c ncsfw.Context) error {
			_ = a.initContextWithApiKey(c)
			_ = a.initContextWithToken(c)
			return next(c)
		})
	}
}

// initContextWithToken cookieからセッションIDを取得して検証
func (a *AuthService) initContextWithToken(c ncsfw.Context) error {
	r := c.Request()
	ctx := r.Context()
	session := a.session
	s, err := session.Get(r, OauthSessionName)
	if err != nil {
		c.Redirect(http.StatusFound, a.Cfg.RootURL.String()+"/login")
		return ncsfw.ErrRedirect
	}
	var loginUser model.LoginUser
	loginUserJson, ok := s.Values[loginKey].(string)
	if !ok {
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "session error", nil)
	}
	if err := loginUser.Deserialize(loginUserJson); err != nil {
		a.log.ErrorWithContext(ctx, err, "login user deserialize")
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "login user deserialize", err)
	}

	if v, err := a.ApiService.GetLoginUser(ctx, loginUser.Email); v == nil {
		return response.ErrJsonWithStatus(err.Code, err.Message, err.Error)
	} else {
		loginUser = *v
	}
	c.SignedIn(true)
	c.SetLoginUser(&loginUser)
	metadata.SetMetadata(c.Request(), metadata.XApiKey, loginUser.Email)
	return nil
}

// initContextWithApiKey Authorization header or basic authorizeからapi keyを取得して検証
func (a *AuthService) initContextWithApiKey(c ncsfw.Context) error {
	header := c.GetHeader("Authorization")
	parts := strings.SplitN(header, " ", 2)
	var keyString string
	if len(parts) == 2 && parts[0] == "Bearer" {
		keyString = parts[1]
	} else {
		username, password, err := utils.DecodeBasicAuthHeader(header)
		if err == nil && username == "api_key" {
			keyString = password
		}
	}

	if keyString == "" {
		return ncsfw.ErrIgnore
	}
	r := c.Request()
	ctx := r.Context()
	loginUser, err := a.ApiService.Decode(ctx, keyString)
	if err != nil {
		return response.ErrJsonWithStatus(err.Code, err.Message, err.Error)
	}
	c.SetRequest(r.WithContext(ctx))
	c.SignedIn(true)
	c.SetLoginUser(loginUser)
	metadata.SetMetadata(c.Request(), metadata.XApiKey, loginUser.Email)
	return nil
}

func (a *AuthService) RoleCheck(role model.RoleLevel) ncsfw.MiddlewareFunc {
	return func(next ncsfw.HandlerFunc) ncsfw.HandlerFunc {
		return func(c ncsfw.Context) error {
			loginUser, ok := c.GetLoginUser().(*model.LoginUser)
			if !ok {
				return response.ErrJsonWithStatus(http.StatusUnauthorized, "", nil)
			}
			level, err := model.String2RoleLevel(loginUser.Role)
			if err != nil {
				return response.ErrJsonWithStatus(http.StatusUnauthorized, "", nil)
			}
			if !role.IsLevelEnabled(level) {
				return response.ErrJsonWithStatus(http.StatusForbidden, "", nil)
			}
			return next(c)
		}
	}
}

func NewAuthMiddleware(cfg *config.Config, socialService social.Service, apiService service.APIKey) *AuthMiddleware {
	return &AuthMiddleware{
		Cfg:           cfg,
		SocialService: socialService,
		ApiService:    apiService,
		log:           logger.New("auth middleware"),
	}
}

func (a *AuthMiddleware) RestMiddleware(session sessions.Store) *AuthService {
	return &AuthService{
		Cfg:           a.Cfg,
		session:       session,
		SocialService: a.SocialService,
		ApiService:    a.ApiService,
		log:           a.log.WithSubModule("auth service"),
	}
}
