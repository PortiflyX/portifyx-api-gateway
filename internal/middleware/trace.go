package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/labstack/echo/v4"
	"github.com/portifyx/portifyx-goutils/pkg/logx"
)

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "0000000000000000"
	}
	return hex.EncodeToString(bytes)
}

// TraceMiddleware ensures every incoming request has trace-id and span-id in context and response headers.
func TraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			traceID := req.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = randomHex(16)
			}
			spanID := req.Header.Get("X-Span-ID")
			if spanID == "" {
				spanID = randomHex(8)
			}

			// Add to response header
			c.Response().Header().Set("X-Trace-ID", traceID)
			c.Response().Header().Set("X-Span-ID", spanID)

			// Propagate down context
			ctx := logx.ContextWithTrace(req.Context(), traceID, spanID)
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}
