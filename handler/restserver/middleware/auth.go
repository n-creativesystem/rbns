package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/internal/apikeygen"
	"github.com/n-creativesystem/rbns/internal/utils"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/tracer"
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
	log           logger.Logger
}

type AuthService struct {
	Cfg           *config.Config
	SocialService social.Service
	session       sessions.Store
	log           logger.Logger
}

func (a *AuthService) middleware(c *contexts.Context, fn func(c *contexts.Context) error) error {
	ctx, span := tracer.Start(c.Request.Context(), "Auth - Middleware")
	defer span.End()
	*c.Request = *c.Request.WithContext(ctx)
	return fn(c)
}

func (a *AuthService) Required(c *contexts.Context) error {
	return a.middleware(c, func(c *contexts.Context) error {
		returnErr := func(err error) bool {
			return err != nil && !(errors.Is(err, contexts.ErrIgnore) || errors.Is(err, contexts.ErrRedirect))
		}

		if err := a.initContextWithApiKey(c); returnErr(err) {
			return err
		}

		if err := a.initContextWithToken(c); returnErr(err) {
			return err
		}
		if c.IsSignedIn {
			c.Next()
		} else {
			return response.ErrJsonWithStatus(http.StatusUnauthorized, "no login", nil)
		}
		return nil
	})
}

func (a *AuthService) NotRequired(c *contexts.Context) error {
	return a.middleware(c, func(c *contexts.Context) error {
		_ = a.initContextWithApiKey(c)
		_ = a.initContextWithToken(c)
		c.Next()
		return nil
	})
}

func (a *AuthService) initContextWithApiKey(c *contexts.Context) error {
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
		return contexts.ErrIgnore
	}
	ctx := c.Request.Context()
	loginUser, res := a.initContextWithAPIKey(ctx, keyString)
	c.Request = c.Request.WithContext(ctx)
	if res.IsError() {
		return res
	}
	c.IsSignedIn = true
	c.LoginUser = loginUser
	return nil
}

func (a *AuthService) initContextWithToken(c *contexts.Context) error {
	session := a.session
	s, err := session.Get(c.Request, OauthSessionName)
	if err != nil {
		c.Redirect(http.StatusFound, a.Cfg.RootURL.String()+"/login")
		return contexts.ErrRedirect
	}
	var loginUser model.LoginUser
	loginUserJson, ok := s.Values[loginKey].(string)
	if !ok {
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "session error", nil)
	}
	if err := loginUser.Deserialize(loginUserJson); err != nil {
		a.log.ErrorWithContext(c, err, "login user deserialize")
		return response.ErrJsonWithStatus(http.StatusUnauthorized, "login user deserialize", err)
	}
	c.IsSignedIn = true
	c.LoginUser = &loginUser
	return nil
}

func (a *AuthService) initContextWithAPIKey(ctx context.Context, keyString string) (*model.LoginUser, response.ErrorResponse) {
	ctx, span := tracer.Start(ctx, "initContextWithApiKey")
	defer span.End()
	decoded, err := apikeygen.Decode(keyString)
	if err != nil {
		return nil, response.ErrJsonWithStatus(http.StatusUnauthorized, err.Error(), nil)
	}

	queryAPIKey := model.GetAPIKeyByNameQuery{
		KeyName: decoded.Name,
	}

	if err := bus.Dispatch(ctx, &queryAPIKey); err != nil {
		return nil, response.ErrJsonWithStatus(http.StatusUnauthorized, err.Error(), nil)
	}

	apiKey := queryAPIKey.Result

	isValid, err := apikeygen.IsValid(decoded, apiKey.Key)
	if err != nil {
		a.log.ErrorWithContext(ctx, err, "Validating API key failed", "API key name", apiKey.Name)
		return nil, response.ErrJsonWithStatus(http.StatusInternalServerError, "Validating API key failed", err)
	}
	if !isValid {
		a.log.ErrorWithContext(ctx, InValidAPIKey, InValidAPIKey.Error(), "API key name", apiKey.Name)
		return nil, response.ErrJsonWithStatus(http.StatusInternalServerError, InValidAPIKey.Error(), nil)
	}

	return a.getLoginUser(ctx, apiKey.ServiceAccountID)
}

func (a *AuthService) getLoginUser(ctx context.Context, id string) (*model.LoginUser, response.ErrorResponse) {
	queryLoginUser := model.GetLoginUserByIDQuery{
		ID: id,
	}
	if err := bus.Dispatch(ctx, &queryLoginUser); err != nil {
		a.log.ErrorWithContext(ctx, err, "Failed to link API key to service account in", "id", queryLoginUser.ID)
		return nil, response.ErrJsonWithStatus(http.StatusUnauthorized, "Failed to link API key to service account in", err)
	}
	return queryLoginUser.Result, response.ErrorResponse{}
}

func NewAuthMiddleware(cfg *config.Config, socialService social.Service) *AuthMiddleware {
	return &AuthMiddleware{
		Cfg:           cfg,
		SocialService: socialService,
		log:           logger.New("auth middleware"),
		// authService:   authService,
	}
}

func (a *AuthMiddleware) RestMiddleware(session sessions.Store) *AuthService {
	return &AuthService{
		Cfg:           a.Cfg,
		session:       session,
		SocialService: a.SocialService,
		log:           a.log.WithSubModule("auth service"),
	}
}
