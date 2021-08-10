package restserver

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/sirupsen/logrus"
)

var (
	pHandler permissionHandle
	rHandler roleHandle
	oHandler organizationHandle
	uHandler userHandle
)

type webUI struct {
	enabled bool
	prefix  string
	root    string
	indexes bool
	baseURL string
}

type Option func(conf *config)

type config struct {
	whiteList []*net.IPNet
	ui        webUI
	debug     bool
}

func WithDebug(conf *config) {
	conf.debug = true
}

func WithWhiteList(whitelistIp string) Option {
	ips := strings.Split(whitelistIp, ",")
	ipNets := make([]*net.IPNet, len(ips))
	for idx, ip := range ips {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			panic(err)
		}
		ipNets[idx] = ipNet
	}
	return func(conf *config) {
		copy(conf.whiteList, ipNets)
	}
}

func WithUI(enabled bool, prefix, root string, indexes bool, baseURL string) Option {
	return func(conf *config) {
		conf.ui = webUI{
			enabled: enabled,
			prefix:  prefix,
			root:    root,
			indexes: indexes,
			baseURL: baseURL,
		}
	}
}

// New is *gin.Engine
//
// Endpoints base is `/api/v1/`
//
// • /permissions
//
// • /roles
//
// • /organizations
//
// • /users
func New(opts ...Option) *gin.Engine {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}
	loggerOpts := []logger.HandlerLogOption{}
	log := logger.NewHandlerLogger()
	if conf.debug {
		log.SetLevel(logrus.DebugLevel)
		loggerOpts = append(loggerOpts, logger.WithGinDebug(logrus.DebugLevel))
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	}
	gin.DefaultWriter = log
	r := gin.New()
	r.Use(logger.RestLogger(loggerOpts...), gin.Recovery())
	if conf.ui.enabled {
		log.Infof("webUI enabled")
		type settings struct {
			BaseURL string `json:"base_url"`
		}
		var s settings
		s.BaseURL = conf.ui.baseURL
		r.GET("/", func(c *gin.Context) {
			c.File("static/index.html")
		})
		r.GET("/index.html", func(c *gin.Context) {
			c.File("static/index.html")
		})
		r.Static("/static", "./static")
		r.GET("/settings.json", func(c *gin.Context) {
			c.JSON(http.StatusOK, &s)
		})
	}
	if len(conf.whiteList) > 0 {
		r.Use(ipFilter(conf.whiteList))
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": http.StatusText(http.StatusNotFound),
		})
	})
	r.Use(errorMiddleware)
	apiV1 := r.Group("/api/v1")
	di.MustInvoke(func(h permissionHandle) {
		pHandler = h
	})
	di.MustInvoke(func(h roleHandle) {
		rHandler = h
	})
	di.MustInvoke(func(h organizationHandle) {
		oHandler = h
	})
	di.MustInvoke(func(h userHandle) {
		uHandler = h
	})
	p := apiV1.Group("/permissions")
	{
		p.GET("", pHandler.findAll)
		p.POST("", pHandler.create)
		p.GET("/:id", pHandler.findById)
		p.PUT("/:id", pHandler.update)
		p.DELETE("/:id", pHandler.delete)
	}
	role := apiV1.Group("/roles")
	{
		role.GET("", rHandler.findAll)
		role.POST("", rHandler.create)
		role.GET("/:id", rHandler.findById)
		role.PUT("/:id", rHandler.update)
		role.DELETE("/:id", rHandler.delete)
		role.GET("/:id/permissions", rHandler.getPermissions)
		role.PUT("/:id/permissions", rHandler.addPermissions)
		role.DELETE("/:id/permissions/:permissionId", rHandler.deletePermissions)
	}
	organization := apiV1.Group("/organizations")
	{
		organization.GET("", oHandler.findAll)
		organization.POST("", oHandler.create)
		organization.GET("/:id", oHandler.findById)
		organization.PUT("/:id", oHandler.update)
		organization.DELETE("/:id", oHandler.delete)
		organization.POST("/:id/users", uHandler.create)
		organization.GET("/:id/users/:key", uHandler.findByKey)
		organization.DELETE("/:id/users/:key", uHandler.delete)
		organization.PUT("/:id/users/:key/roles", uHandler.addRole)
		organization.DELETE("/:id/users/:key/roles/:roleId", uHandler.deleteRole)
	}
	return r
}
