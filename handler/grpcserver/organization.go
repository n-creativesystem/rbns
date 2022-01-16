package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type organizationServer struct {
	orgSvc         service.Organization
	orgAggregation service.OrganizationAggregation
}

var _ protobuf.OrganizationServer = (*organizationServer)(nil)

func NewOrganizationService(orgSvc service.Organization, orgAggregation service.OrganizationAggregation) protobuf.OrganizationServer {
	return &organizationServer{orgSvc: orgSvc, orgAggregation: orgAggregation}
}

// Organization
func (s *organizationServer) Create(ctx context.Context, in *protobuf.OrganizationEntity) (*protobuf.OrganizationEntity, error) {
	org, err := s.orgSvc.Create(ctx, in.GetName(), in.GetDescription())
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := &protobuf.OrganizationEntity{
		Id:          org.ID.String(),
		Name:        org.Name,
		Description: org.Description,
	}
	return out, nil
}

func (s *organizationServer) FindById(ctx context.Context, in *protobuf.OrganizationKey) (*protobuf.OrganizationEntity, error) {
	organization, err := s.orgSvc.FindById(ctx, in.GetId())
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoconv.NewOrganizationEntityByModel(*organization), nil
}

func (s *organizationServer) FindAll(ctx context.Context, in *emptypb.Empty) (*protobuf.OrganizationEntities, error) {
	organizations, err := s.orgSvc.FindAll(ctx)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	protoOrganizations := make([]*protobuf.OrganizationEntity, len(organizations))
	for idx, organization := range organizations {
		protoOrganizations[idx] = protoconv.NewOrganizationEntityByModel(organization)
	}
	out := &protobuf.OrganizationEntities{
		Organizations: protoOrganizations,
	}
	return out, nil
}

func (s *organizationServer) Update(ctx context.Context, in *protobuf.OrganizationUpdateEntity) (*emptypb.Empty, error) {
	if err := s.orgSvc.Update(ctx, in.GetId(), in.GetName(), in.GetDescription()); err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) Delete(ctx context.Context, in *protobuf.OrganizationKey) (*emptypb.Empty, error) {
	if err := s.orgSvc.Delete(ctx, in.GetId()); err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) AddUser(ctx context.Context, in *protobuf.AddOrganizationUser) (*emptypb.Empty, error) {
	if err := s.orgAggregation.AddUsers(ctx, in.GetId(), []string{in.GetUserId()}); err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) DeleteUser(ctx context.Context, in *protobuf.DeleteOrganizationUser) (*emptypb.Empty, error) {
	if err := s.orgAggregation.DeleteUsers(ctx, in.GetUserId(), []string{in.GetUserId()}); err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
