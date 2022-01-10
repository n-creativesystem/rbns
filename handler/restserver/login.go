package restserver

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
	"github.com/n-creativesystem/rbns/handler/restserver/response"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"golang.org/x/oauth2"
)

const (
	loginKey             = "loginUser"
	OauthSessionName     = "oauth_session"
	OauthStateCookieName = "oauth_state"
	OauthPKCECookieName  = "oauth_code_verifier"
)

func (s *HTTPServer) GetOAuthProvider(c *contexts.Context) error {
	result := make(map[string]interface{})
	mp := s.socialService.GetOAuthInfoProviders()
	for key, value := range mp {
		if value.Enabled {
			result[key] = map[string]interface{}{
				"name": value.Name,
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
	return nil
}

func (s *HTTPServer) OAuthLogin(c *contexts.Context) error {
	r := c.Request
	session, err := s.store.Get(r, OauthSessionName)
	if err != nil {
		return err
	}
	name := c.Params.ByName("name")
	provider := s.socialService.GetOAuthInfoProvider(name)
	connector, err := s.socialService.GetConnector(name)
	if err != nil {
		s.log.ErrorWithContext(c, err, "social service connector")
		s.handleOAuthLoginError(c, response.ErrJson("social service connector", err))
		return nil
	}
	errorParam := c.Query("error")
	if errorParam != "" {
		errorDesc := errors.New(c.Query("error_description"))
		s.log.ErrorWithContext(c, errorDesc, errorParam)
		s.handleOAuthLoginError(c, response.ErrJson(errorParam, errorDesc))
		return nil
	}

	code := c.Query("code")
	if code == "" {
		opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOnline}
		if provider.UsePKCE {
			ascii, pkce, err := GenerateCodeChallenge()
			if err != nil {
				s.log.ErrorWithContext(c, err, "Generating PKCE failed")
				s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "An internal error occurred", err))
				return nil
			}
			session.Values[OauthPKCECookieName] = ascii
			// cookies.WriteSessionCookie(c.Writer, OauthPKCECookieName, ascii, s.Cfg.OAuthCookieMaxAge)

			opts = append(opts,
				oauth2.SetAuthURLParam("code_challenge", pkce),
				oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			)
		}
		state, err := GenStateString()
		if err != nil {
			s.log.ErrorWithContext(c, err, "Generating state string failed")
			s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "An internal error occurred", err))
			return nil
		}
		hashedState := hashStatecode(state, s.Cfg.SecretKey, provider.ClientSecret)
		session.Values[OauthStateCookieName] = hashedState
		if err := session.Save(r, c.Writer); err != nil {
			s.log.ErrorWithContext(c, err, "session save error")
			s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "session save error", err))
		}
		if provider.HostedDomain != "" {
			opts = append(opts, oauth2.SetAuthURLParam("hd", provider.HostedDomain))
		}
		c.Redirect(http.StatusFound, connector.AuthCodeURL(state, opts...))
		return nil
	}
	defer func() {
		if err := session.Save(r, c.Writer); err != nil {
			s.log.ErrorWithContext(c, err, "session save error")
			s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "session save error", err))
		}
	}()

	cookieState, ok := session.Values[OauthStateCookieName].(string)
	if !ok {
		s.log.ErrorWithContext(c, err, "login.OAuthLogin(missing saved state)")
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "login.OAuthLogin(missing saved state)", err))
		return nil
	}
	delete(session.Values, OauthStateCookieName)
	queryState := hashStatecode(c.Query("state"), s.Cfg.SecretKey, provider.ClientSecret)
	if cookieState != queryState {
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "login.OAuthLogin(state mismatch)", nil))
		return nil
	}

	oauthClient, err := s.socialService.GetOAuthHttpClient(name)
	if err != nil {
		s.log.ErrorWithContext(c, err, "Failed to create OAuth http client")
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "login.OAuthLogin("+err.Error()+")", err))
		return nil
	}

	oauthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, oauthClient)
	opts := []oauth2.AuthCodeOption{}

	if codeVerifier, ok := session.Values[OauthPKCECookieName].(string); ok {
		delete(session.Values, OauthPKCECookieName)
		if codeVerifier != "" {
			opts = append(opts,
				oauth2.SetAuthURLParam("code_verifier", codeVerifier),
			)
		}
	}

	// get token from provider
	token, err := connector.Exchange(oauthCtx, code, opts...)
	if err != nil {
		s.log.ErrorWithContext(c, err, "login.OAuthLogin(NewTransportWithCode)")
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, "login.OAuthLogin(NewTransportWithCode)", err))
		return nil
	}
	token.TokenType = "Bearer"

	// set up oauth2 client
	client := connector.Client(oauthCtx, token)

	// get user info
	userInfo, err := connector.UserInfo(client, token)
	if err != nil {
		s.log.ErrorWithContext(c, err, fmt.Sprintf("login.OAuthLogin(get info from %s)", name))
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusInternalServerError, fmt.Sprintf("login.OAuthLogin(get info from %s)", name), err))
		return nil
	}

	// validate that we got at least an email address
	if userInfo.Email == "" {
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusUnauthorized, "login provider didn't return an email address", nil))
		return nil
	}

	// validate that the email is allowed to login to grafana
	if !connector.IsEmailAllowed(userInfo.Email) {
		s.handleOAuthLoginError(c, response.ErrJsonWithStatus(http.StatusUnauthorized, "required email domain not fulfilled", nil))
		return nil
	}

	if userInfo.Role == "" {
		userInfo.Role = string(model.ROLE_VIEWER)
	}

	loginUser := buildExternalUserInfo(token, userInfo, name)
	session.Values[loginKey] = loginUser.Serialize()
	if err := syncUser(c.Request.Context(), loginUser, connector); err != nil {
		s.log.ErrorWithContext(c, err, "database save error")
		u := path.Join(s.Cfg.RootURL.Path, "/logout")
		s.Redirect(c, u)
	}
	// metrics.MApiLoginOAuth.Inc()
	if redirectTo, ok := session.Values["redirect_to"].(string); ok {
		if redirectTo, err := url.QueryUnescape(redirectTo); err == nil && len(redirectTo) > 0 {
			if err := s.ValidateRedirectTo(redirectTo); err == nil {
				delete(session.Values, "redirect_to")
				s.Redirect(c, redirectTo)
				return nil
			}
			s.log.DebugWithContext(c, "Ignored invalid redirect_to cookie value", "redirect_to", redirectTo)
		}
	}
	s.Redirect(c, s.Cfg.RootURL.String())
	return nil
}

func (s *HTTPServer) handleOAuthLoginError(c *contexts.Context, err response.ErrorResponse) {
	s.log.ErrorWithContext(c, err, "login")
	c.Abort()
	u := path.Join(s.Cfg.RootURL.Path, "/logint")
	s.Redirect(c, u)
}

func (s *HTTPServer) Logout(c *contexts.Context) error {
	r := c.Request
	session, err := s.store.Get(r, OauthSessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	_ = session.Save(r, c.Writer)
	s.log.InfoWithContext(c, "Successful Logout")
	s.Redirect(c, s.Cfg.RootURL.String())
	return nil
}

func buildExternalUserInfo(token *oauth2.Token, userInfo *social.BasicUserInfo, name string) *model.LoginUser {
	user := (&model.LoginUser{
		ID:        userInfo.Id,
		UseName:   userInfo.Name,
		Email:     userInfo.Email,
		Role:      userInfo.Role,
		Groups:    userInfo.Groups,
		OAuthToke: token,
	}).SetOAuthName(name)
	return user
}

func syncUser(ctx context.Context, user *model.LoginUser, connector social.SocialConnector) error {
	cmd := &model.UpsertLoginUserCommand{
		User:          user,
		SignupAllowed: connector.IsSignupAllowed(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func GenStateString() (string, error) {
	rnd := make([]byte, 32)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rnd), nil
}

func hashStatecode(code, secretKey, seed string) string {
	hashBytes := sha256.Sum256([]byte(code + secretKey + seed))
	return hex.EncodeToString(hashBytes[:])
}

func GenerateCodeVerifier() (codeVerifier []byte, err error) {
	raw := make([]byte, 96)
	_, err = rand.Read(raw)
	if err != nil {
		return nil, err
	}
	ascii := make([]byte, 128)
	base64.RawURLEncoding.Encode(ascii, raw)
	return ascii, nil
}

func GenerateCodeChallenge() (string, string, error) {
	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])
	return string(codeVerifier), codeChallenge, nil
}
