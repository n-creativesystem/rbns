package client

import (
	"context"
	"math"
	"time"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GRPCClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn

	permissionClient   protobuf.PermissionClient
	roleClient         protobuf.RoleClient
	organizationClient protobuf.OrganizationClient
	userClient         protobuf.UserClient
}

func NewGRPCClient(grpcAddress string) (*GRPCClient, error) {
	return NewGRPCClientWithContext(grpcAddress, context.Background())
}

func NewGRPCClientWithContext(grpcAddress string, baseCtx context.Context) (*GRPCClient, error) {
	return NewGRPCClientWithContextTLS(grpcAddress, baseCtx, "", "")
}

func NewGRPCClientWithContextTLS(grpcAddress string, baseCtx context.Context, certFile, commonName string) (*GRPCClient, error) {
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
	return &GRPCClient{
		ctx:                ctx,
		cancel:             cancel,
		conn:               conn,
		permissionClient:   protobuf.NewPermissionClient(conn),
		roleClient:         protobuf.NewRoleClient(conn),
		organizationClient: protobuf.NewOrganizationClient(conn),
		userClient:         protobuf.NewUserClient(conn),
	}, nil
}

func (c *GRPCClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return c.ctx.Err()
}

func (c *GRPCClient) Target() string {
	return c.conn.Target()
}

func (c *GRPCClient) Permissions() Permissions {
	return &permissionClient{
		ctx:    c.ctx,
		client: c.permissionClient,
	}
}

func (c *GRPCClient) Roles() Roles {
	return &roleClient{
		ctx:    c.ctx,
		client: c.roleClient,
	}
}

func (c *GRPCClient) Organizations() Organizations {
	return &organizationClient{
		ctx:    c.ctx,
		client: c.organizationClient,
	}
}

func (c *GRPCClient) Users() Users {
	return &userClient{
		ctx:    c.ctx,
		client: c.userClient,
	}
}
