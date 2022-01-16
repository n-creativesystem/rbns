package config

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
)

var (
	RootURL                = ""
	CookieSecure           = true
	CookieSameSiteDisabled = false
	CookieSameSiteMode     = http.SameSiteStrictMode
)

type Config struct {
	ConfigFile  string
	GrpcPort    int
	GatewayPort int
	StorageType string `validate:"required"`
	// RootURL is http(s)://domain/subpath
	RootURL *url.URL
	// SubPath サブパス リバプロに対応するため
	SubPath           string
	KeyPairs          []string `validate:"required"`
	LogoutURL         string
	MetadataURL       string
	CertificateFile   string
	KeyFile           string
	CommonName        string
	StaticFilePath    string
	MultiTenant       bool
	SecretKey         string
	OAuthCookieMaxAge time.Duration

	CacheDriverName string
	CacheDataSource string

	LoginCookieName string

	OAuthRaw *Section

	DatabaseRaw *Section
}

func NewFlags2Config(flagSet *pflag.FlagSet) (*Config, error) {
	flags := flags{flagSet}
	conf := Config{}

	conf.ConfigFile = flags.MustString("configFile")
	conf.GrpcPort = flags.MustInt("grpcPort")
	conf.GatewayPort = flags.MustInt("httpPort")

	conf.StorageType = flags.MustString("storageType")

	// /subpath
	conf.SubPath = flags.MustString("subPath")
	rootPath := flags.MustString("rootUrl")
	rootURL, _ := url.Parse(rootPath)
	rootURL.Path = flags.MustString("subPath")
	// http://localhost:xxxx/subpath
	conf.RootURL = rootURL
	conf.StaticFilePath = flags.MustString("staticFilePath")
	conf.KeyPairs = flags.MustStringArray("keyPairs")
	oauthCookieMaxAge := flags.MustInt("oauth_cookie_max_age")
	conf.OAuthCookieMaxAge = time.Duration(oauthCookieMaxAge) * time.Minute
	conf.SecretKey = flags.MustString("hash_secret_key")
	conf.loadRaw()

	validate := validator.New()
	err := validate.Struct(conf)
	return &conf, err
}

func (c *Config) loadRaw() {
	c.OAuthRaw = &Section{key: "auth"}
	c.DatabaseRaw = &Section{key: "database"}
}

// func (c *Config) IsSAML() bool {
// 	return len(c.MetadataURL) > 0
// }

// func (c *Config) IsOIDC() bool {
// 	return len(c.IssuerURL) > 0
// }
