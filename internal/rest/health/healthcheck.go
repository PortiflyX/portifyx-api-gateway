package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /healthcheck (Liveness probe)
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Readiness handles GET /healthcheck-readiness (Readiness probe)
func (h *HealthHandler) Readiness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ready",
		"service": "portifyx-api-gateway",
	})
}
