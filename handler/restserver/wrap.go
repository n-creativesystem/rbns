package restserver

import (
	"net/http"
	"strings"

	"github.com/n-creativesystem/rbns/handler/metadata"
	"github.com/n-creativesystem/rbns/ncsfw"
)

func (hs *HTTPServer) WrapGateway(h http.Handler) ncsfw.HandlerFunc {
	return func(c ncsfw.Context) error {
		r := c.Request()
		r.URL.Path = strings.Replace(r.URL.Path, "/api/v1/g", "/api/v1", 1)

		metadata.SetMetadata(r, metadata.XTenantID, c.GetTenant())
		hs.gateway.ServeHTTP(c.Writer(), r)
		return nil
	}
}
