package auth

// import (
// 	"encoding/json"
// 	"net"
// 	"net/http"
// 	"net/url"
// 	"time"

// 	"github.com/crewjam/saml"
// 	"github.com/crewjam/saml/samlsp"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/sessions"
// 	"github.com/n-creativesystem/rbns/config"
// 	"github.com/sirupsen/logrus"
// )

// type samlImpl struct {
// 	*samlsp.Middleware
// }

// func (s *samlImpl) Use(r *gin.Engine) {
// 	r.Any("saml/*action", gin.WrapH(s.Middleware))
// 	r.Use(s.middleware)
// }

// func (s *samlImpl) middleware(c *gin.Context) {
// 	w := c.Writer
// 	r := c.Request
// 	session, err := s.Middleware.Session.GetSession(r)
// 	if err != nil {
// 		logrus.Errorf("Session error: %s", err)
// 	}
// 	if err == samlsp.ErrNoSession {
// 		s.Middleware.HandleStartAuthFlow(w, r)
// 		c.Abort()
// 		return
// 	}
// 	if session != nil {
// 		c.Request = r.WithContext(samlsp.ContextWithSession(r.Context(), session))
// 		buf, err := json.Marshal(session)
// 		if err != nil {
// 			logrus.Errorf("Session Json Marshal: %s", err)
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
// 			return
// 		}
// 		c.Set(AuthKey, string(buf))
// 		c.Next()
// 		return
// 	}
// 	logrus.Error(err)
// 	c.AbortWithStatus(http.StatusUnauthorized)
// }

// type samlConfig struct {
// 	MetadataURL string
// 	RootURL     string
// 	IDPLogout   string
// }

// func NewSAML(conf *config.Config, store sessions.Store) (Middleware, error) {
// 	samlConf := samlConfig{
// 		MetadataURL: conf.MetadataURL,
// 		RootURL:     conf.RootURL.String(),
// 		IDPLogout:   conf.LogoutURL,
// 	}
// 	sp, err := serviceProvider(samlConf, store)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &samlImpl{Middleware: sp}, nil
// }

// func serviceProvider(conf samlConfig, store sessions.Store) (*samlsp.Middleware, error) {
// 	idpMetadataURL, err := url.Parse(conf.MetadataURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rootURL, err := url.Parse(conf.RootURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opts := samlsp.Options{
// 		URL:            *rootURL,
// 		IDPMetadataURL: idpMetadataURL,
// 	}
// 	samlSP, err := samlsp.New(opts)
// 	samlSP.Session = newSessionProvider(opts, store)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return samlSP, nil
// }

// const name = "token"
// const valueKey = "token"

// type samlStore struct {
// 	store    sessions.Store
// 	name     string
// 	domain   string
// 	httpOnly bool
// 	secure   bool
// 	maxAge   time.Duration
// 	codec    samlsp.SessionCodec
// }

// var _ samlsp.SessionProvider = (*samlStore)(nil)

// func newSessionProvider(opts samlsp.Options, store sessions.Store) samlsp.SessionProvider {
// 	cookieName := name
// 	if opts.CookieName != "" {
// 		cookieName = opts.CookieName
// 	}

// 	maxAge := time.Hour
// 	if opts.CookieMaxAge > 0 {
// 		maxAge = opts.CookieMaxAge
// 	}
// 	// for backwards compatibility, support CookieDomain
// 	cookieDomain := opts.URL.Host
// 	if opts.CookieDomain != "" {
// 		cookieDomain = opts.CookieDomain
// 	}

// 	// for backwards compatibility, support CookieSecure
// 	cookieSecure := opts.URL.Scheme == "https"
// 	if opts.CookieSecure {
// 		cookieSecure = true
// 	}
// 	return &samlStore{
// 		store:    store,
// 		name:     cookieName,
// 		domain:   cookieDomain,
// 		maxAge:   maxAge,
// 		httpOnly: true,
// 		secure:   cookieSecure,
// 		codec:    DefaultSessionCodec(opts),
// 	}
// }

// func (s *samlStore) CreateSession(w http.ResponseWriter, r *http.Request, assertion *saml.Assertion) error {
// 	if domain, _, err := net.SplitHostPort(s.domain); err == nil {
// 		s.domain = domain
// 	}

// 	samlSession, err := s.codec.New(assertion)
// 	if err != nil {
// 		return err
// 	}

// 	value, err := s.codec.Encode(samlSession)
// 	if err != nil {
// 		return err
// 	}

// 	session, err := s.store.New(r, s.name)
// 	if err != nil {
// 		return err
// 	}

// 	session.Values[valueKey] = value

// 	session.Options.Domain = s.domain
// 	session.Options.MaxAge = int(s.maxAge.Seconds())
// 	session.Options.HttpOnly = s.httpOnly
// 	session.Options.Secure = s.secure || r.URL.Scheme == "https"
// 	session.Options.Path = "/"

// 	return session.Save(r, w)
// }

// func (s *samlStore) DeleteSession(w http.ResponseWriter, r *http.Request) error {
// 	session, err := s.store.Get(r, s.name)
// 	if err != nil {
// 		return err
// 	}

// 	session.Options.MaxAge = -1

// 	return session.Save(r, w)
// }

// func (s *samlStore) GetSession(r *http.Request) (samlsp.Session, error) {
// 	session, err := s.store.Get(r, s.name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	v, ok := session.Values[valueKey]
// 	if !ok {
// 		return nil, samlsp.ErrNoSession
// 	}

// 	samlSession, err := s.codec.Decode(v.(string))
// 	if err != nil {
// 		return nil, samlsp.ErrNoSession
// 	}
// 	return samlSession, nil
// }
