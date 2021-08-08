package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/handler/restserver/internal/file"
)

func staticServe(urlPrefix string, fs file.ServeFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
