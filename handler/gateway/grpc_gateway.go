package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/handler/gateway/marshaler"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/protobuf"
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
		XTenantID:  f(XTenantID, "fake"),
		XIndexKey:  f(XIndexKey, "index"),
		XRequestID: req.Header.Get(XRequestID),
	})
}

type GRPCGateway struct {
	mux http.Handler
}

func New(conf *config.Config) (*GRPCGateway, error) {
	grpcAddress, certFile, commonName := fmt.Sprintf(":%d", conf.GrpcPort), conf.CertificateFile, conf.KeyFile
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
			logger.ErrorWithContext(ctx, err, "failed to register handler from endpoint")
			return nil, err
		}
	}

	return &GRPCGateway{
		mux: mux,
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
		// protobuf.RegisterResourceHandlerFromEndpoint,
	}
}
