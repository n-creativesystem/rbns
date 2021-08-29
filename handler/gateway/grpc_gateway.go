package gateway

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/n-creativesystem/rbns/handler/gateway/marshaler"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func responseFilter(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	switch resp.(type) {
	default:
		w.Header().Set("Content-Type", marshaler.DefaultContentType)
	}

	return nil
}

const (
	XTenantID  = "X-Tenant-ID"
	XIndexKey  = "X-Index-Key"
	XRequestID = "X-Request-ID"
)

func withMetadata(ctx context.Context, req *http.Request) metadata.MD {
	f := func(key, default_ string) string {
		if v := req.Header.Get(key); v != "" {
			return v
		}
		return default_
	}
	return metadata.New(map[string]string{
		XTenantID:  f(XTenantID, "default"),
		XIndexKey:  f(XIndexKey, "index"),
		XRequestID: req.Header.Get(XRequestID),
	})
}

type GRPCGateway struct {
	grpcAddress string
	mux         http.Handler
	logger      *logrus.Logger
}

type config struct {
	saml              bool
	saml_payload_name string
	role              string
}

func New(grpcAddress, certFile, commonName, baseURL string, log *logrus.Logger) (*GRPCGateway, error) {
	if baseURL == "" {
		baseURL = "/"
	}
	var err error
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(math.MaxInt64),
			grpc.MaxCallRecvMsgSize(math.MaxInt64),
		),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                1 * time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
	}

	baseCtx := context.TODO()
	ctx, cancel := context.WithCancel(baseCtx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	mux := runtime.NewServeMux(
		runtime.WithMetadata(withMetadata),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, new(marshaler.GatewayMarshaler)),
		runtime.WithForwardResponseOption(responseFilter),
	)
	if certFile == "" {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		var creds credentials.TransportCredentials
		creds, err = credentials.NewClientTLSFromFile(certFile, commonName)
		if err != nil {
			return nil, err
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}
	endpoints := endpoints()
	for _, endpoint := range endpoints {
		if err := endpoint(ctx, mux, grpcAddress, dialOpts); err != nil {
			log.Errorf("failed to register handler from endpoint: %s", err.Error())
			return nil, err
		}
	}
	type settings struct {
		BaseURL string `json:"base_url"`
	}
	var s settings
	s.BaseURL = baseURL
	gin.DefaultWriter = logger.NewWriter(log)
	r := gin.New()
	r.Use(logger.RestLogger(log), gin.Recovery())
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})
	r.GET("/settings.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, &s)
	})
	r.Any("/api/*gateway", gin.WrapH(mux))

	return &GRPCGateway{
		grpcAddress: grpcAddress,
		mux:         r,
		logger:      log,
	}, nil
}

func (g *GRPCGateway) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.mux.ServeHTTP(w, req)
}

var Swagger io.Reader = bytes.NewReader(protobuf.Swagger)

type endpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func endpoints() []endpoint {
	return []endpoint{
		protobuf.RegisterPermissionHandlerFromEndpoint,
		protobuf.RegisterRoleHandlerFromEndpoint,
		protobuf.RegisterUserHandlerFromEndpoint,
		protobuf.RegisterOrganizationHandlerFromEndpoint,
	}
}
