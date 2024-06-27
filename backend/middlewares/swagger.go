package middlewares

import (
	"hbd/docs"
	"strings"

	"github.com/gin-gonic/gin"
)

func SwaggerHostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			docs.SwaggerInfo.Host = c.Request.Host
		}
		c.Next()
	}
}
