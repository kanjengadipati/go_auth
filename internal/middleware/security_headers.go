package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	apiContentSecurityPolicy  = "default-src 'none'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'"
	docsContentSecurityPolicy = "default-src 'self'; img-src 'self' data: https:; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Writer.Header()
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("Referrer-Policy", "no-referrer")
		headers.Set("X-Permitted-Cross-Domain-Policies", "none")
		headers.Set("Content-Security-Policy", contentSecurityPolicy(c.Request.URL.Path))
		if requestUsesHTTPS(c) {
			headers.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

func contentSecurityPolicy(path string) string {
	if path == "/docs" || strings.HasPrefix(path, "/docs/") {
		return docsContentSecurityPolicy
	}
	return apiContentSecurityPolicy
}

func requestUsesHTTPS(c *gin.Context) bool {
	if c.Request != nil && c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
