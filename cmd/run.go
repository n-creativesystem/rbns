package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/handler/gateway"
	"github.com/n-creativesystem/rbns/handler/grpcserver"
	"github.com/n-creativesystem/rbns/infra"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type webUI struct {
	enable  bool
	prefix  string
	root    string
	indexes bool
	baseURL string
}

type databaseConfig struct {
	dialector    string
	masterDSN    string
	slaveDSN     string
	maxIdleConns int
	maxOpenConns int
	maxLifeTime  int
	dbType       string
}

type server struct {
	enabled bool
	port    int
}

var (
	signals = []os.Signal{
		os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	}
	httpSrv     server
	gRPCSrv     server
	storageType string
	debug       bool
	whiteList   string
	secure      bool
	ui          webUI
	database    databaseConfig
	omitHeaders string
	apiKey      string
	runCmd      = &cobra.Command{
		Use:   "run",
		Short: "application start",
		Long:  "application start",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logrus.New()
			log.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
			})
			log.SetLevel(logrus.DebugLevel)
			factory, settings, err := infra.NewFactory(storageType)
			if err != nil {
				return err
			}
			defer factory.Close()
			if err := factory.Initialize(settings, log); err != nil {
				return err
			}
			reader := factory.Reader()
			writer := factory.Writer()
			if err := run(context.Background(), reader, writer, log); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	cobra.OnInitialize(initialize(nil, nil))
	allLogLevel := make([]string, len(logrus.AllLevels))
	for i, level := range logrus.AllLevels {
		allLogLevel[i] = level.String()
	}
	runCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "config file. if omitted, rbns.yaml in /etc and home directory will be searched")
	runCmd.PersistentFlags().IntVar(&httpSrv.port, "gateway-port", 8080, "http port")
	runCmd.PersistentFlags().IntVar(&gRPCSrv.port, "grpc-port", 8888, "grpc port")
	runCmd.PersistentFlags().StringVar(&storageType, "storage-type", "postgres", "persistent storage type")
	runCmd.PersistentFlags().StringVar(&ui.baseURL, "baseURL", "/", "base url")
}

func run(ctx context.Context, reader repository.Reader, writer repository.Writer, log *logrus.Logger) error {
	var (
		eg                     *errgroup.Group
		grpcLister, httpLister net.Listener
		err                    error
	)
	defer func() {
		if grpcLister != nil {
			grpcLister.Close()
		}
		if httpLister != nil {
			httpLister.Close()
		}
	}()

	eg, ctx = errgroup.WithContext(ctx)

	grpcAddr := fmt.Sprintf(":%d", gRPCSrv.port)
	grpcLister, err = net.Listen("tcp", grpcAddr)
	if err != nil {
		logrus.Fatalln(err)
	}
	eg.Go(func() error {
		logrus.Printf("GRPC Server: %s", grpcAddr)
		return runGrpc(ctx, grpcLister, reader, writer, log)
	})

	httpAddr := fmt.Sprintf(":%d", httpSrv.port)
	httpLister, err = net.Listen("tcp", httpAddr)
	if err != nil {
		logrus.Fatalln(err)
	}
	eg.Go(func() error {
		logrus.Printf("REST Server: %s", httpAddr)
		return runRest(ctx, httpLister, log)
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
			return nil
		}
		return err
	}
	return nil
}

func runRest(ctx context.Context, li net.Listener, log *logrus.Logger) error {
	// opts := []restserver.Option{
	// 	restserver.WithGRPC(fmt.Sprintf(":%d", gRPCSrv.port), "", ""),
	// }
	// if debug {
	// 	opts = append(opts, restserver.WithDebug)
	// }
	// if whiteList != "" {
	// 	opts = append(opts, restserver.WithWhiteList(whiteList))
	// }
	// opts = append(opts, restserver.WithUI(ui.enable, ui.prefix, ui.root, ui.indexes, ui.baseURL))
	restSrv, err := gateway.New(fmt.Sprintf(":%d", gRPCSrv.port), "", "", ui.baseURL, log)
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		Handler:      restSrv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errCh := make(chan error)
	go func() {
		if err := httpServer.Serve(li); err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		cancelCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		return httpServer.Shutdown(cancelCtx)
	case err := <-errCh:
		return err
	}
}

func runGrpc(ctx context.Context, li net.Listener, reader repository.Reader, writer repository.Writer, log *logrus.Logger) error {
	opts := []grpcserver.Option{}
	if secure {
		opts = append(opts, grpcserver.WithSecure)
	}
	grpcServer := grpcserver.New(reader, writer, log, opts...)
	errCh := make(chan error)
	go func() {
		if err := grpcServer.Serve(li); err != nil {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}
