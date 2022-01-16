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
	handler_metadata "github.com/n-creativesystem/rbns/handler/metadata"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"github.com/n-creativesystem/rbns/protobuf"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/proto"
)

func responseFilter(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	switch resp.(type) {
	default:
		w.Header().Set("Content-Type", marshaler.DefaultContentType)
	}

	return nil
}

type GRPCGateway struct {
	mux http.Handler
}

func New(conf *config.Config) (*GRPCGateway, error) {
	grpcAddress, certFile, commonName := fmt.Sprintf(":%d", conf.GrpcPort), conf.CertificateFile, conf.KeyFile
	var err error
	otelgrpcOpts := []otelgrpc.Option{
		otelgrpc.WithTracerProvider(tracer.GetTracerProvider()),
		otelgrpc.WithPropagators(tracer.GetPropagation()),
	}
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
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(otelgrpcOpts...)),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(otelgrpcOpts...)),
	}

	baseCtx := context.Background()
	ctx, cancel := context.WithCancel(baseCtx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	mux := runtime.NewServeMux(
		runtime.WithMetadata(handler_metadata.WithMetadata),
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
		protobuf.RegisterUserHandlerFromEndpoint,
		protobuf.RegisterOrganizationHandlerFromEndpoint,
		// userHandlerFromEndpoint,
		// protobuf.RegisterResourceHandlerFromEndpoint,
	}
}
