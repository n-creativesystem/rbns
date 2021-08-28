package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type organizationServer struct {
	*proto.UnimplementedOrganizationServer
	svc service.OrganizationService
}

var _ proto.OrganizationServer = (*organizationServer)(nil)

func newOrganizationService(svc service.OrganizationService) proto.OrganizationServer {
	return &organizationServer{svc: svc}
}

// Organization
func (s *organizationServer) Create(ctx context.Context, in *proto.OrganizationEntity) (*proto.OrganizationEntity, error) {
	org, err := s.svc.Create(ctx, in.GetName(), in.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := &proto.OrganizationEntity{
		Id:          *org.GetID(),
		Name:        *org.GetName(),
		Description: org.GetDescription(),
	}
	return out, nil
}

func (s *organizationServer) FindById(ctx context.Context, in *proto.OrganizationKey) (*proto.OrganizationEntity, error) {
	organization, err := s.svc.FindById(ctx, in.GetId())
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoconv.NewOrganizationEntityByModel(*organization), nil
}

func (s *organizationServer) FindAll(ctx context.Context, in *emptypb.Empty) (*proto.OrganizationEntities, error) {
	organizations, err := s.svc.FindAll(ctx)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	protoOrganizations := make([]*proto.OrganizationEntity, len(organizations))
	for idx, organization := range organizations {
		protoOrganizations[idx] = protoconv.NewOrganizationEntityByModel(organization)
	}
	out := &proto.OrganizationEntities{
		Organizations: protoOrganizations,
	}
	return out, nil
}

func (s *organizationServer) Update(ctx context.Context, in *proto.OrganizationUpdateEntity) (*emptypb.Empty, error) {
	if err := s.svc.Update(ctx, in.GetId(), in.GetName(), in.GetDescription()); err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) Delete(ctx context.Context, in *proto.OrganizationKey) (*emptypb.Empty, error) {
	if err := s.svc.Delete(ctx, in.GetId()); err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
