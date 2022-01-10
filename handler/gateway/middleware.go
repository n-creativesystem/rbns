package gateway

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/logger"
)

func parseSAML(headerName string, claimNames []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger(c)
		payload := c.GetHeader(headerName) // "x-saml-payload"
		buf, err := base64.RawURLEncoding.DecodeString(payload)
		if err != nil {
			log.Error(err, "")
		}
		mp := map[string]interface{}{}
		if err := json.Unmarshal(buf, &mp); err != nil {
			log.Error(err, "")
		}
		m := mp
		for _, key := range claimNames {
			if v, ok := m[key].(map[string]interface{}); ok {
				m = v
				continue
			}
			// if v, ok := m[key]; ok {
			// 	log.Info(v)
			// }
		}
	}
}
