package grpcserver

import (
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type authorizationServer interface {
	auth.AuthorizationServer
}

type authorizationServerImpl struct {
	// svc service.Resource
}

// var _ auth.AuthorizationServer = (*authorizationServerImpl)(nil)

// func (a *authorizationServerImpl) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
// 	var unauthorizedStatus = &auth.CheckResponse{
// 		Status: &rpcstatus.Status{
// 			Code: int32(rpc.UNAUTHENTICATED),
// 		},
// 		HttpResponse: &auth.CheckResponse_DeniedResponse{
// 			DeniedResponse: &auth.DeniedHttpResponse{
// 				Status: &envoy_type.HttpStatus{
// 					Code: envoy_type.StatusCode_Unauthorized,
// 				},
// 				Body: "Authorization Header malformed or not provided",
// 			},
// 		},
// 	}
// 	var permissionDenied = &auth.CheckResponse{
// 		Status: &rpcstatus.Status{
// 			Code: int32(rpc.PERMISSION_DENIED),
// 		},
// 		HttpResponse: &auth.CheckResponse_DeniedResponse{
// 			DeniedResponse: &auth.DeniedHttpResponse{
// 				Status: &envoy_type.HttpStatus{
// 					Code: envoy_type.StatusCode_Unauthorized,
// 				},
// 				Body: "PERMISSION_DENIED",
// 			},
// 		},
// 	}
// 	log.Println(">>> Authorization called check()")

// 	b, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
// 	if err == nil {
// 		log.Println("Inbound Headers: ")
// 		log.Println((string(b)))
// 	}

// 	ct, err := json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
// 	if err == nil {
// 		log.Println("Context Extensions: ")
// 		log.Println((string(ct)))
// 	}
// 	jwtPayloadbase64 := req.Attributes.Request.Http.Headers["x-jwt-payload"]
// 	jwtPayloadbase64Buf, err := base64.RawURLEncoding.DecodeString(jwtPayloadbase64)
// 	if err != nil {
// 		return unauthorizedStatus, nil
// 	}
// 	jwtMap := map[string]interface{}{}
// 	err = json.Unmarshal(jwtPayloadbase64Buf, &jwtMap)
// 	if err != nil {
// 		return unauthorizedStatus, nil
// 	}
// 	userkey, ok := jwtMap["sub"].(string)
// 	if !ok {
// 		return permissionDenied, nil
// 	}
// 	organization := "default"
// 	// ok = a.svc.Authorized(ctx, req.Attributes.Request.Http.Method, req.Attributes.Request.Http.Path, organization, userkey)
// 	if ok {
// 		return &auth.CheckResponse{
// 			Status: &rpcstatus.Status{
// 				Code: int32(rpc.OK),
// 			},
// 			HttpResponse: &auth.CheckResponse_OkResponse{
// 				OkResponse: &auth.OkHttpResponse{
// 					Headers: []*core.HeaderValueOption{
// 						{
// 							Header: &core.HeaderValue{
// 								Key:   "x-custom-header-from-authz",
// 								Value: "some value",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}, nil
// 	} else {
// 		return permissionDenied, nil
// 	}
// }

// func newAuthz(svc service.Resource) authorizationServer {
// 	return &authorizationServerImpl{
// 		svc: svc,
// 	}
// }

// func envoyAuthzRegister(s *grpc.Server, srv authorizationServer) {
// 	auth.RegisterAuthorizationServer(s, srv)
// }
