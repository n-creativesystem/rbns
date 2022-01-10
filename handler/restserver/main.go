package restserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/handler/gateway"
	"github.com/n-creativesystem/rbns/handler/restserver/middleware"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/storage"
)

func New(
	gateway *gateway.GRPCGateway,
	conf *config.Config,
	authMiddleware *middleware.AuthMiddleware,
	store sessions.Store,
	socialService social.Service,
	_ storage.Factory,
) (*http.Server, error) {
	server := &HTTPServer{
		log:            logger.New("rest server"),
		gateway:        gateway,
		socialService:  socialService,
		authMiddleware: authMiddleware,
		store:          store,
		Cfg:            conf,
	}
	server.registerRouting()
	httpAddr := fmt.Sprintf(":%d", conf.GatewayPort)
	s := &http.Server{
		Handler:      server,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         httpAddr,
	}
	return s, nil
}
