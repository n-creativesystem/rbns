package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"github.com/n-creativesystem/rbns/version"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "application start",
		Long:  "application start",
		Run:   run,
		PreRun: func(cmd *cobra.Command, args []string) {
			initialize(cmd)
		},
	}
	flags := cmd.PersistentFlags()
	flags.String("configFile", "", "config file. if omitted, rbns.yaml in /etc and home directory will be searched")
	flags.Int("httpPort", 8080, "http port")
	flags.Int("grpcPort", 8888, "grpc port")
	flags.String("storageType", "postgres", "persistent storage type")

	flags.String("rootUrl", "/", "base url")
	flags.String("subPath", "/", "sub path")
	flags.String("staticFilePath", "static", "static file path")
	flags.StringArray("keyPairs", []string{"secure"}, "secure cookie secret keys")
	flags.String("logoutUrl", "/", "logout url")
	flags.String("hash_secret_key", "secret", "generate state hash secret key")
	flags.Int("oauth_cookie_max_age", 60, "oauth cookie time(m)")
	flags.String("impl_name", "sql", "bus implement name")

	return cmd
}

type servers struct {
	conf       *config.Config
	restServer *http.Server
	grpcServer *grpc.Server
}

func newServer(restServer *http.Server, grpcServer *grpc.Server, conf *config.Config) *servers {
	return &servers{
		conf:       conf,
		restServer: restServer,
		grpcServer: grpcServer,
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, _ = tracer.InitOpenTelemetryWithService(ctx, "role based N security", tracer.Service{
		Name:    "rbns",
		Version: version.Version,
	})
	flags := cmd.PersistentFlags()
	s, err := initializeRun(flags)
	if err != nil {
		logger.FatalWithContext(ctx, err, "initializeRun")
	}
	var (
		grpcLister, httpLister net.Listener
	)

	grpcAddr := fmt.Sprintf(":%d", s.conf.GrpcPort)
	grpcLister, err = net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.FatalWithContext(ctx, err, "grpc listener error")
	}
	httpAddr := fmt.Sprintf(":%d", s.conf.GatewayPort)
	httpLister, err = net.Listen("tcp", httpAddr)
	if err != nil {
		logger.FatalWithContext(ctx, err, "rest listener error")
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		logger.InfoWithContext(ctx, "grpc sever start")
		_ = s.grpcServer.Serve(grpcLister)
		return nil
	})
	eg.Go(func() error {
		logger.InfoWithContext(ctx, "rest sever start")
		_ = s.restServer.Serve(httpLister)
		return nil
	})
	eg.Go(func() error {
		return signal(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	if err := eg.Wait(); err != nil {
		if err == SignalReceived {
			return
		}
		logger.FatalWithContext(ctx, err, "error group wait error")
	}

	cancelCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := s.restServer.Shutdown(cancelCtx); err != nil {
		logger.ErrorWithContext(cancelCtx, err, "rest server shutdown")
	}
	s.grpcServer.GracefulStop()
}
