package main

import (
	"flag"
	"fmt"
)

type items []string

func (i *items) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *items) Set(value string) error {
	*i = append(*i, value)
	return nil
}

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
}

type server struct {
	enabled bool
	port    int
}

var (
	httpSrv   server
	gRPCSrv   server
	debug     bool
	whiteList string
	secure    bool
	ui        webUI
	envFiles  items
	database  databaseConfig
)

func init() {
	flag.BoolVar(&httpSrv.enabled, "http", true, "start rest api server")
	flag.IntVar(&httpSrv.port, "httpPort", 8080, "http port")
	flag.BoolVar(&gRPCSrv.enabled, "grpc", true, "start gRPC server")
	flag.IntVar(&gRPCSrv.port, "grpcPort", 8888, "grpc port")
	flag.BoolVar(&debug, "debug", true, "debug mode")
	flag.StringVar(&whiteList, "whitelist", "", "ip address whitelist(CIDR)")
	flag.BoolVar(&secure, "secure", true, "api key guard")
	flag.BoolVar(&ui.enable, "ui", false, "setting is web ui")
	flag.StringVar(&ui.prefix, "uiPrefix", "/", "static file prefix")
	flag.Var(&envFiles, "env", ".env file name")
	flag.StringVar(&database.dialector, "dialector", "postgres", "database driver dialector")
	flag.StringVar(&database.masterDSN, "masterDSN", "", "master database source name")
	flag.StringVar(&database.slaveDSN, "slaveDSN", "", "slave database source name")
	flag.IntVar(&database.maxIdleConns, "maxIdleConns", 10, "database max idle connections")
	flag.IntVar(&database.maxOpenConns, "maxOpenConns", 100, "database max open connections")
	flag.IntVar(&database.maxLifeTime, "maxLifeTime", 1, "database connection life time")
	flag.StringVar(&ui.root, "staticRoot", "./static", "web ui static root")
	flag.BoolVar(&ui.indexes, "staticIndexes", false, "web ui static allow indexes")
	flag.StringVar(&ui.baseURL, "baseURL", "/", "base url")
}
