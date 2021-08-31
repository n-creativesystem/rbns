module github.com/n-creativesystem/rbns/client

go 1.16

replace github.com/n-creativesystem/rbns/protobuf v0.0.0 => ../protobuf

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/n-creativesystem/rbns/protobuf v0.0.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
