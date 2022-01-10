package middleware

// type ApiKey interface {
// 	UnaryServerInterceptor() grpc.UnaryServerInterceptor
// 	StreamServerInterceptor() grpc.StreamServerInterceptor
// }

// func NewApiKey(svc service.ApiKey) ApiKey {
// 	return &apiKey{
// 		svc: svc,
// 	}
// }

// type apiKey struct {
// 	err error
// }

// func (api *apiKey) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		switch {
// 		case api.authorization(ctx):
// 		}
// 		if api.err != nil {
// 			err = status.Error(codes.Unauthenticated, api.err.Error())
// 			return
// 		}
// 		return handler(ctx, req)
// 	}
// }

// func (api *apiKey) StreamServerInterceptor() grpc.StreamServerInterceptor {
// 	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		ctx := ss.Context()
// 		switch {
// 		case api.authorization(ctx):
// 		}
// 		if api.err != nil {
// 			return status.Error(codes.Unauthenticated, api.err.Error())
// 		}
// 		return handler(srv, ss)
// 	}
// }

// const headerAuthorize = "authorization"

// func (api *apiKey) authorization(ctx context.Context) bool {
// 	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
// 	if val == "" {
// 		return false
// 	}
// 	parts := strings.SplitN(val, " ", 2)
// 	var keyString string
// 	if len(parts) == 2 && parts[0] == "Bearer" {
// 		keyString = parts[1]
// 	} else {
// 		username, password, err := utils.DecodeBasicAuthHeader(val)
// 		if err == nil && username == "api_key" {
// 			keyString = password
// 		}
// 	}
// 	if keyString == "" {
// 		return false
// 	}
// 	apiKeyGen, err := model.DecodeApiKey(keyString)
// 	if err != nil {
// 		api.err = err
// 		return true
// 	}
// 	key, err := api.svc.Check(ctx, apiKeyGen.Name)
// 	if err != nil {
// 		api.err = err
// 		return true
// 	}
// 	ok, err := apiKeyGen.IsValid(key)
// 	if err != nil {
// 		api.err = err
// 		return true
// 	}
// 	return ok
// }
