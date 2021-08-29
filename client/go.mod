module github.com/n-creativesystem/rbns/client

go 1.16

replace github.com/n-creativesystem/rbns/protobuf v0.0.0 => ../protobuf

require (
	github.com/n-creativesystem/rbns/protobuf v0.0.0
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
