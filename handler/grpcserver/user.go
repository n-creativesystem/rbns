package grpcserver

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
	"github.com/n-creativesystem/rbns/utilsconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServer struct {
	svc            service.User
	orgAggregation service.OrganizationAggregation
}

var _ protobuf.UserServer = (*userServer)(nil)

func NewUserServer(svc service.User, orgAggregation service.OrganizationAggregation) protobuf.UserServer {
	return &userServer{svc: svc, orgAggregation: orgAggregation}
}

// User
func (s *userServer) Create(ctx context.Context, in *protobuf.UserCreateKey) (*emptypb.Empty, error) {
	err := s.svc.Create(ctx, in.GetId(), in.GetName())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, in *protobuf.UserKey) (*emptypb.Empty, error) {
	err := s.svc.Delete(ctx, in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, err
}

func (s *userServer) FindById(ctx context.Context, in *protobuf.UserKey) (*protobuf.UserEntity, error) {
	u, err := s.svc.FindById(ctx, in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	roles := make([]*protobuf.RoleEntity, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, protoconv.NewRoleEntityByModel(r))
	}
	permissions := make([]*protobuf.PermissionEntity, 0, len(u.Permissions))
	for _, p := range u.Permissions {
		permissions = append(permissions, protoconv.NewPermissionEntityByModel(p))
	}
	out := protobuf.UserEntity{
		Id:          u.ID,
		Name:        u.Name,
		Roles:       roles,
		Permissions: permissions,
	}
	return &out, nil
}

func (s *userServer) AddRoles(ctx context.Context, in *protobuf.UserRole) (*emptypb.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.orgAggregation.AddUserRoles(ctx, in.GetOrganizationId(), in.GetId(), roles)
	if err != nil {
		if err == model.ErrNoDataFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) DeleteRoles(ctx context.Context, in *protobuf.UserRole) (*emptypb.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.orgAggregation.AddUserRoles(ctx, in.GetOrganizationId(), in.GetId(), roles)
	if err != nil {
		if err == model.ErrNoDataFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) DeleteRole(ctx context.Context, in *protobuf.UserRoleDelete) (*emptypb.Empty, error) {
	return s.DeleteRoles(ctx, &protobuf.UserRole{
		Id:             in.GetId(),
		OrganizationId: in.OrganizationId,
		Roles: []*protobuf.RoleKey{
			{
				Id: in.RoleId,
			},
		},
	})
}

func (s *userServer) Uploader(ctx context.Context, in *protobuf.UserUploader) (*emptypb.Empty, error) {
	chunk := in.GetData()
	b64data := chunk[strings.IndexByte(chunk, ',')+1:]
	buf, err := utilsconv.Base64ToByte(b64data)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewBuffer(buf))
	fileType := in.GetFileType()
	for _, opt := range fileType.Option {
		if opt.Key == "lazyQuotes" {
			if len(opt.Values) >= 1 {
				reader.LazyQuotes, _ = strconv.ParseBool(opt.Values[0])
			}
		}
	}

	// 先頭はヘッダーカラムの為読み飛ばし
	if _, err := reader.Read(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		id := utilsconv.CopyString(record[0])
		name := utilsconv.CopyString(record[1])
		if err := s.svc.Create(ctx, id, name); err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}
