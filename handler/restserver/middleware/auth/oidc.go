package auth

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/base64"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"path"
// 	"sync/atomic"
// 	"time"

// 	oidc "github.com/coreos/go-oidc"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gin-gonic/gin/render"
// 	"github.com/gorilla/sessions"
// 	"github.com/n-creativesystem/rbns/config"
// 	"github.com/sirupsen/logrus"
// 	"golang.org/x/oauth2"
// )

// var allowedSigningAlgs = map[string]bool{
// 	oidc.RS256: true,
// 	oidc.RS384: true,
// 	oidc.RS512: true,
// 	oidc.ES256: true,
// 	oidc.ES384: true,
// 	oidc.ES512: true,
// 	oidc.PS256: true,
// 	oidc.PS384: true,
// 	oidc.PS512: true,
// }

// type oidcPath string

// func (o oidcPath) String() string {
// 	return path.Join("/auth", string(o))
// }

// const (
// 	oidcLoginPath    oidcPath = "/login"
// 	oidcCallBackPath oidcPath = "/callback"
// 	oidcLogoutPath   oidcPath = "/logout"
// )

// type oidcImpl struct {
// 	provider *oidc.Provider
// 	verifier atomic.Value
// 	conf     oauth2.Config
// 	store    sessions.Store
// 	rootURL  string
// 	opts     oidcConfig
// }

// func (o *oidcImpl) setVerifier(v *oidc.IDTokenVerifier) {
// 	o.verifier.Store(v)
// }

// func (o *oidcImpl) idTokenVerifier() (*oidc.IDTokenVerifier, bool) {
// 	if v := o.verifier.Load(); v != nil {
// 		return v.(*oidc.IDTokenVerifier), true
// 	}
// 	return nil, false
// }

// func (o *oidcImpl) Use(r *gin.Engine) {
// 	r.GET(oidcLoginPath.String(), o.Login)
// 	r.GET(oidcCallBackPath.String(), o.CallBack)
// 	r.GET(oidcLogoutPath.String(), o.Logout)
// 	r.Use(o.middleware)
// }

// func (o *oidcImpl) middleware(c *gin.Context) {
// 	r := c.Request
// 	session, err := o.store.Get(r, name)
// 	if err != nil {
// 		logrus.Error(c.AbortWithError(http.StatusInternalServerError, err))
// 		return
// 	}
// 	rawIdToken, ok := session.Values["id_token"].(string)
// 	if !ok {
// 		c.Render(-1, o.loginRedirect(c.Request))
// 		c.Abort()
// 		return
// 	}
// 	c.Set(AuthKey, rawIdToken)
// 	c.Next()
// }

// func (o *oidcImpl) Login(c *gin.Context) {
// 	w := c.Writer
// 	r := c.Request
// 	b := make([]byte, 32)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	state := base64.StdEncoding.EncodeToString(b)
// 	session, err := o.store.Get(r, name)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	session.Values["state"] = state
// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, o.conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
// }

// func (o *oidcImpl) CallBack(c *gin.Context) {
// 	w := c.Writer
// 	r := c.Request
// 	session, err := o.store.Get(r, name)
// 	if err != nil {
// 		_ = c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	if r.URL.Query().Get("state") != session.Values["state"] {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid state parameter"})
// 		return
// 	}

// 	token, err := o.conf.Exchange(r.Context(), r.URL.Query().Get("code"))
// 	if err != nil {
// 		logrus.Warnf("no token found: %v", err)
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 		return
// 	}

// 	rawIDToken, ok := token.Extra("id_token").(string)
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "No id_token field in oauth2 token."})
// 		return
// 	}

// 	verifier, ok := o.idTokenVerifier()
// 	if !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "oidc: authenticator not initialized"})
// 		return
// 	}
// 	idToken, err := verifier.Verify(r.Context(), rawIDToken)
// 	if err != nil {
// 		logrus.Error(c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to verify ID Token: "+err.Error())))
// 		return
// 	}
// 	mp := map[string]interface{}{}
// 	if err := idToken.Claims(&mp); err != nil {
// 		logrus.Error(c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("oidc: parse claims: %v", err)))
// 		return
// 	}
// 	buf, _ := json.Marshal(&mp)

// 	delete(session.Values, "state")
// 	session.Values["id_token"] = string(buf)
// 	err = session.Save(r, w)
// 	if err != nil {
// 		logrus.Error(err)
// 		_ = c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.Render(-1, o.indexRedirect(c.Request))
// }

// func (o *oidcImpl) Logout(c *gin.Context) {
// 	w := c.Writer
// 	r := c.Request
// 	session, err := o.store.Get(r, name)
// 	if err != nil {
// 		logrus.Errorf("logout: %s", c.AbortWithError(http.StatusInternalServerError, err))
// 		return
// 	}
// 	session.Options.MaxAge = -1
// 	if err := session.Save(r, w); err != nil {
// 		logrus.Errorf("logout: %s", err)
// 	}
// 	c.Render(-1, o.logoutRedirect(c.Request))
// }

// func (o *oidcImpl) loginRedirect(r *http.Request) render.Render {
// 	url := o.opts.joinSubPath(oidcLoginPath.String())
// 	return render.Redirect{
// 		Code:     http.StatusSeeOther,
// 		Location: url,
// 		Request:  r,
// 	}
// }

// func (o *oidcImpl) indexRedirect(r *http.Request) render.Render {
// 	return render.Redirect{
// 		Code:     http.StatusSeeOther,
// 		Location: o.opts.RootURL.String(),
// 		Request:  r,
// 	}
// }

// func (o *oidcImpl) logoutRedirect(r *http.Request) render.Render {
// 	if o.opts.LogoutURL == "" {
// 		return render.Redirect{
// 			Code:     http.StatusSeeOther,
// 			Location: o.opts.RootURL.String(),
// 			Request:  r,
// 		}
// 	} else {
// 		return render.Redirect{
// 			Code:     http.StatusSeeOther,
// 			Location: o.opts.LogoutURL,
// 			Request:  r,
// 		}
// 	}
// }

// type oidcConfig struct {
// 	*config.Config
// }

// func (o *oidcConfig) joinSubPath(elem ...string) string {
// 	elem = append([]string{o.SubPath}, elem...)
// 	return path.Join(elem...)
// }

// func NewOIDC(config *config.Config, store sessions.Store) (Middleware, error) {
// 	opts := oidcConfig{
// 		Config: config,
// 	}
// 	issuerUrl, err := url.Parse("")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if issuerUrl.Scheme != "https" {
// 		return nil, fmt.Errorf("'oidc-issuer-url' (%q) has invalid scheme (%q), require 'https'", "opts.IssuerURL", issuerUrl.Scheme)
// 	}
// 	rootURL, err := url.Parse(opts.RootURL.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	callbackURL := *rootURL
// 	callbackURL.Path = path.Join(callbackURL.Path, oidcCallBackPath.String())
// 	supportedSigningAlgs := []string{}
// 	if len(supportedSigningAlgs) == 0 {
// 		supportedSigningAlgs = []string{oidc.RS256}
// 	}
// 	for _, alg := range supportedSigningAlgs {
// 		if !allowedSigningAlgs[alg] {
// 			return nil, fmt.Errorf("oidc: unsupported signing alg: %q", alg)
// 		}
// 	}

// 	client := &http.Client{Timeout: 30 * time.Second}

// 	ctx := context.Background()
// 	ctx = oidc.ClientContext(ctx, client)
// 	now := time.Now
// 	verifierConfig := &oidc.Config{
// 		ClientID:             "opts.ClientID",
// 		SupportedSigningAlgs: supportedSigningAlgs,
// 		Now:                  now,
// 	}
// 	provider, err := oidc.NewProvider(ctx, "opts.IssuerURL")
// 	if err != nil {
// 		logrus.Errorf("oidc authenticator: initializing plugin: %v", err)
// 		return nil, err
// 	}
// 	verifier := provider.Verifier(verifierConfig)
// 	conf := oauth2.Config{
// 		ClientID:     "opts.ClientID",
// 		ClientSecret: "opts.ClientSecret",
// 		RedirectURL:  callbackURL.String(),
// 		Endpoint:     provider.Endpoint(),
// 		Scopes:       []string{oidc.ScopeOpenID, "profile"},
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	authenticator := &oidcImpl{
// 		provider: provider,
// 		verifier: atomic.Value{},
// 		conf:     conf,
// 		store:    store,
// 		rootURL:  rootURL.String(),
// 		opts:     opts,
// 	}
// 	authenticator.setVerifier(verifier)
// 	return authenticator, nil
// }
