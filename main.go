package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/n-creativesystem/rbns/infra"

	_ "github.com/n-creativesystem/rbns/service"

	"github.com/joho/godotenv"
	"github.com/n-creativesystem/rbns/handler/grpcserver"
	"github.com/n-creativesystem/rbns/handler/restserver"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/utilsconv"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func setEnvFlag() {
	flag.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(utilsconv.ToSnakeCase(f.Name))
		if s := os.Getenv(strings.ToUpper(name)); s != "" {
			err := f.Value.Set(s)
			if err != nil {
				panic(err)
			}
		}
	})
}

func main() {
	logger.SetFormatter(logrus.StandardLogger())
	dao.Register()
	setEnvFlag()
	flag.Parse()
	for _, envFile := range envFiles {
		_ = godotenv.Load(envFile)
	}
	setEnvFlag()
	dbOpts := []dao.Option{
		dao.WithDialector(database.dialector),
		dao.WithMasterDSN(database.masterDSN),
		dao.WithSlaveDSN(database.slaveDSN),
		dao.WithMaxIdleConn(database.maxIdleConns),
		dao.WithMaxOpenConns(database.maxOpenConns),
		dao.WithMaxLifeTime(database.maxLifeTime),
		dao.WithMigration,
	}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		dbOpts = append(dbOpts, dao.Debug)
	}
	db := dao.New(dbOpts...)
	if err := run(context.Background(), db); err != nil {
		logrus.Fatalln(err)
	}
}

func run(ctx context.Context, db dao.DataBase) error {
	var (
		eg                     *errgroup.Group
		grpcLister, httpLister net.Listener
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

	if gRPCSrv.enabled {
		var err error
		grpcAddr := fmt.Sprintf(":%d", gRPCSrv.port)
		grpcLister, err = net.Listen("tcp", grpcAddr)
		defer grpcLister.Close()
		if err != nil {
			logrus.Fatalln(err)
		}
		eg.Go(func() error {
			logrus.Printf("GRPC Server: %s", grpcAddr)
			return runGrpc(ctx, db, grpcLister)
		})
	}
	if httpSrv.enabled {
		var err error
		httpAddr := fmt.Sprintf(":%d", httpSrv.port)
		httpLister, err = net.Listen("tcp", httpAddr)
		if err != nil {
			logrus.Fatalln(err)
		}
		eg.Go(func() error {
			logrus.Printf("REST Server: %s", httpAddr)
			return runRest(ctx, httpLister, db)
		})
	}
	eg.Go(func() error {
		return signal(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	return eg.Wait()
}

func runRest(ctx context.Context, li net.Listener, db dao.DataBase) error {
	opts := []restserver.Option{}
	if debug {
		opts = append(opts, restserver.WithDebug)
	}
	if whiteList != "" {
		opts = append(opts, restserver.WithWhiteList(whiteList))
	}
	opts = append(opts, restserver.WithUI(ui.enable, ui.prefix, ui.root, ui.indexes, ui.baseURL))
	restSrv := restserver.New(opts...)
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

func runGrpc(ctx context.Context, db dao.DataBase, li net.Listener) error {
	opts := []grpcserver.Option{}
	if secure {
		opts = append(opts, grpcserver.WithSecure)
	}
	grpcServer := grpcserver.New(db, opts...)
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
