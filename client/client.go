package client

import (
	"context"
	"io"
	"math"
	"time"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GRPCClient interface {
	Target() string
}

type RBNS interface {
	Permissions(ctx context.Context) Permissions
	Roles(ctx context.Context) Roles
	Organizations(ctx context.Context) Organizations
	Users(ctx context.Context) Users
	Resource(ctx context.Context) Resource
	io.Closer
	GRPCClient
}

type clientImpl struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn

	permissionClient   protobuf.PermissionClient
	roleClient         protobuf.RoleClient
	organizationClient protobuf.OrganizationClient
	userClient         protobuf.UserClient
	resourceClient     protobuf.ResourceClient
}

func New(grpcAddress string) (RBNS, error) {
	return NewWithContext(grpcAddress, context.Background())
}

func NewWithContext(grpcAddress string, baseCtx context.Context) (RBNS, error) {
	return NewWithContextTLS(grpcAddress, baseCtx, "", "")
}

func NewWithContextTLS(grpcAddress string, baseCtx context.Context, certFile, commonName string) (RBNS, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(math.MaxInt32),
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	if certFile == "" {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(certFile, commonName)
		if err != nil {
			return nil, err
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}
	ctx, cancel := context.WithCancel(baseCtx)
	conn, err := grpc.DialContext(ctx, grpcAddress, dialOpts...)
	if err != nil {
		cancel()
		return nil, err
	}
	return &clientImpl{
		ctx:                ctx,
		cancel:             cancel,
		conn:               conn,
		permissionClient:   protobuf.NewPermissionClient(conn),
		roleClient:         protobuf.NewRoleClient(conn),
		organizationClient: protobuf.NewOrganizationClient(conn),
		userClient:         protobuf.NewUserClient(conn),
		resourceClient:     protobuf.NewResourceClient(conn),
	}, nil
}

func (c *clientImpl) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.ctx.Err()
}

func (c *clientImpl) Target() string {
	return c.conn.Target()
}

func (c *clientImpl) Permissions(ctx context.Context) Permissions {
	return &permissionClient{
		ctx:    ctx,
		client: c.permissionClient,
	}
}

func (c *clientImpl) Roles(ctx context.Context) Roles {
	return &roleClient{
		ctx:    ctx,
		client: c.roleClient,
	}
}

func (c *clientImpl) Organizations(ctx context.Context) Organizations {
	return &organizationClient{
		ctx:    ctx,
		client: c.organizationClient,
	}
}

func (c *clientImpl) Users(ctx context.Context) Users {
	return &userClient{
		ctx:    ctx,
		client: c.userClient,
	}
}

func (c *clientImpl) Resource(ctx context.Context) Resource {
	return &resource{
		ctx:    ctx,
		client: c.resourceClient,
	}
}
