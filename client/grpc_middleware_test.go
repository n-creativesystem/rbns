package client_test

import (
	"context"
	"net"
	"testing"

	"github.com/n-creativesystem/rbns/client"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop"
	pb "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func serve(sOpt ...grpc.ServerOption) *grpc.Server {
	lis = bufconn.Listen(1024 * 1024)
	s := grpc.NewServer(sOpt...)
	// TestServiceServerを登録する
	pb.RegisterTestServiceServer(s, interop.NewTestServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return s
}

func TestMiddleware(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	c, err := client.New("localhost:8888")
	if !assert.NoError(t, err) {
		return
	}
	mp := client.MethodPermissions{
		"/grpc.testing.TestService/EmptyCall": []string{"create:test"},
	}
	fn := func(ctx context.Context) (newCtx context.Context, userKey string, organizationName string, err error) {
		return ctx, "test", "default", nil
	}
	s := serve(
		grpc.StreamInterceptor(client.StreamServerInterceptor(c, mp, fn)),
		grpc.UnaryInterceptor(client.UnaryServerInterceptor(c, mp, fn)),
	)
	defer s.Stop()
	ctx := context.Background()
	// ダイアル関数
	dial := func(context.Context, string) (net.Conn, error) {
		// lisはグローバルに宣言されている変数
		return lis.Dial()
	}
	// gPRCコネクションの生成
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		append([]grpc.DialOption{
			grpc.WithContextDialer(dial),
			grpc.WithInsecure(),
		})...,
	)
	if err != nil {
		t.Fatalf("fialed to dial: %v", err)
	}
	defer conn.Close()

	// テスト用gRPCクライアントの生成
	client := pb.NewTestServiceClient(conn)
	// gRPCのメソッド呼び出し
	interop.DoEmptyUnaryCall(client)
}
