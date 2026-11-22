package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/portifyx/portifyx-api-gateway/config"
	"github.com/portifyx/portifyx-api-gateway/internal/rest"
	"github.com/stretchr/testify/assert"
)

func TestGatewayHealthEndpoints(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			AppName: "portifyx-api-gateway",
			AppPort: "8080",
		},
	}

	srv := rest.NewEchoServer(cfg)
	e := srv.GetEcho()

	// 1. Test GET /healthcheck
	req1 := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rec1 := httptest.NewRecorder()
	e.ServeHTTP(rec1, req1)

	assert.Equal(t, http.StatusOK, rec1.Code)
	var body1 map[string]string
	assert.NoError(t, json.Unmarshal(rec1.Body.Bytes(), &body1))
	assert.Equal(t, "ok", body1["status"])

	// 2. Test GET /healthcheck-readiness
	req2 := httptest.NewRequest(http.MethodGet, "/healthcheck-readiness", nil)
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)

	assert.Equal(t, http.StatusOK, rec2.Code)
	var body2 map[string]interface{}
	assert.NoError(t, json.Unmarshal(rec2.Body.Bytes(), &body2))
	assert.Equal(t, "ready", body2["status"])
}
