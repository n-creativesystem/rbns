package grpcserver

// type resourceServer struct {
// 	*protobuf.UnimplementedResourceServer
// 	svc service.Resource
// }

// var _ protobuf.ResourceServer = (*resourceServer)(nil)

// func NewResourceServer(svc service.Resource) protobuf.ResourceServer {
// 	return &resourceServer{
// 		svc: svc,
// 	}
// }

// func (s *resourceServer) Find(ctx context.Context, in *protobuf.ResourceKey) (*protobuf.ResourceResponse, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	result, err := s.svc.Find(ctx, id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
// 	res := protoconv.NewResourceByModel(*result)
// 	return res, nil
// }

// func (s *resourceServer) Exists(ctx context.Context, in *protobuf.ResourceKey) (*protobuf.ResourceExistsResponse, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	isExists, _ := s.svc.Exists(ctx, id)
// 	return &protobuf.ResourceExistsResponse{IsExists: isExists}, nil
// }

// func (s *resourceServer) Delete(ctx context.Context, in *protobuf.ResourceKey) (*emptypb.Empty, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := s.svc.Delete(ctx, id); err != nil {
// 		return nil, err
// 	}
// 	return &emptypb.Empty{}, nil
// }

// func (s *resourceServer) Update(ctx context.Context, in *protobuf.ResourceSaveRequest) (*emptypb.Empty, error) {
// 	return s.Create(ctx, in)
// }

// func (s *resourceServer) Create(ctx context.Context, in *protobuf.ResourceSaveRequest) (*emptypb.Empty, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	empty := &emptypb.Empty{}
// 	if err := s.svc.Save(ctx, id, in.Description, in.PermissionNames...); err != nil {
// 		return empty, status.Error(codes.Internal, err.Error())
// 	}
// 	return empty, nil
// }

// func (s *resourceServer) FindAll(ctx context.Context, _ *emptypb.Empty) (*protobuf.ResourceResponses, error) {
// 	results, err := s.svc.FindAll(ctx)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
// 	resp := &protobuf.ResourceResponses{
// 		Resources: make([]*protobuf.ResourceResponse, 0, len(results)),
// 	}
// 	for _, result := range results {
// 		res := protoconv.NewResourceByModel(*result)
// 		resp.Resources = append(resp.Resources, res)
// 	}
// 	return resp, nil
// }

// func (s *resourceServer) Migration(ctx context.Context, in *protobuf.ResourceSaveRequest) (*emptypb.Empty, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := s.svc.Migration(ctx, id, in.Description, in.PermissionNames...); err != nil {
// 		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
// 	}
// 	return &emptypb.Empty{}, nil
// }

// func (s *resourceServer) GetPermissions(ctx context.Context, in *protobuf.ResourceKey) (*protobuf.PermissionEntities, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	permissions, err := s.svc.GetPermissions(ctx, id)
// 	if err != nil {
// 		if err == model.ErrNoData {
// 			return nil, status.Error(codes.NotFound, err.Error())
// 		}
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
// 	res := protobuf.PermissionEntities{
// 		Permissions: make([]*protobuf.PermissionEntity, len(permissions)),
// 	}
// 	for idx, permission := range permissions {
// 		res.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
// 	}
// 	return &res, nil
// }

// func (s *resourceServer) AddPermissions(ctx context.Context, in *protobuf.ResourceReleationPermissions) (*emptypb.Empty, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	permissionIds := make([]string, len(in.GetPermissions()))
// 	if len(permissionIds) == 0 {
// 		return &emptypb.Empty{}, nil
// 	}
// 	for idx, permission := range in.GetPermissions() {
// 		permissionIds[idx] = permission.GetId()
// 	}
// 	if err := s.svc.AddPermissions(ctx, id, permissionIds); err != nil {
// 		if err == model.ErrNoData {
// 			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
// 		}
// 		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
// 	}
// 	return &emptypb.Empty{}, nil
// }

// func (s *resourceServer) DeletePermission(ctx context.Context, in *protobuf.ResourceReleationPermission) (*emptypb.Empty, error) {
// 	return s.DeletePermissions(ctx, &protobuf.ResourceReleationPermissions{
// 		Id: in.Id,
// 		Permissions: []*protobuf.PermissionKey{
// 			{
// 				Id: in.PermissionId,
// 			},
// 		},
// 	})
// }

// func (s *resourceServer) DeletePermissions(ctx context.Context, in *protobuf.ResourceReleationPermissions) (*emptypb.Empty, error) {
// 	id, err := model.NewKey(in.GetId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	permissionIds := make([]string, len(in.GetPermissions()))
// 	if len(permissionIds) == 0 {
// 		return &emptypb.Empty{}, nil
// 	}
// 	for idx, permission := range in.GetPermissions() {
// 		permissionIds[idx] = permission.GetId()
// 	}
// 	if err := s.svc.DeletePermissions(ctx, id, permissionIds); err != nil {
// 		if err == model.ErrNoData {
// 			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
// 		}
// 		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
// 	}
// 	return &emptypb.Empty{}, nil
// }
