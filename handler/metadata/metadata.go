package metadata

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

const (
	XTenantID  = "X-Tenant-ID"
	XRequestID = "X-Request-ID"
	XApiKey    = "X-Api-Key"
)

func getHeader(req *http.Request, key, default_ string) string {
	if v := req.Header.Get(key); v != "" {
		return v
	}
	return default_
}

func WithMetadata(ctx context.Context, req *http.Request) metadata.MD {
	return metadata.New(map[string]string{
		XTenantID:  getHeader(req, XTenantID, ""),
		XRequestID: req.Header.Get(XRequestID),
		XApiKey:    req.Header.Get(XApiKey),
	})
}

func SetMetadata(req *http.Request, key, value string) {
	req.Header.Add(key, value)
}
