package restserver

import (
	"github.com/n-creativesystem/rbns/handler/restserver/contexts"
)

func (s *HTTPServer) gatewayWrap(c *contexts.Context) error {
	s.gateway.ServeHTTP(c.Writer, c.Request)
	return nil
}
