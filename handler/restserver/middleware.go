package restserver

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ipFilter(ipNets []*net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := net.ParseIP(c.ClientIP())
		for _, ip := range ipNets {
			if ip.Contains(clientIP) {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func setApiKey(localApiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("Authorization")
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
		if strings.EqualFold(localApiKey, apiKey) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
