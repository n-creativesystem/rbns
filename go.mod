module github.com/n-creativesystem/rbns

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/envoyproxy/go-control-plane v0.9.10-0.20210907150352-cf90f659a021
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/google/wire v0.5.0
	github.com/gorilla/sessions v1.2.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jackc/pgconn v1.10.1
	github.com/jackc/pgx/v4 v4.14.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0
	github.com/lyft/protoc-gen-star v0.6.0 // indirect
	github.com/n-creativesystem/rbns/protobuf v0.0.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.1.7
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.7
	github.com/wader/gormstore/v2 v2.0.0
	github.com/xhit/go-str2duration/v2 v2.0.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/contrib/propagators/b3 v1.3.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/exporters/prometheus v0.26.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/sdk/export/metric v0.26.0
	go.opentelemetry.io/otel/sdk/metric v0.26.0
	go.opentelemetry.io/otel/trace v1.3.0
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.2.2
	gorm.io/driver/postgres v1.2.3
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.4
	gorm.io/plugin/dbresolver v1.1.0
)

replace github.com/n-creativesystem/rbns/protobuf v0.0.0 => ./protobuf

replace github.com/sirupsen/logrus v1.8.1 => github.com/n-creativesystem/logrus v1.9.1
