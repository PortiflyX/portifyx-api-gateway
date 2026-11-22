package rest

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/portifyx/portifyx-api-gateway/config"
	"github.com/portifyx/portifyx-api-gateway/internal/middleware"
	"github.com/portifyx/portifyx-api-gateway/internal/rest/health"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

func NewEchoServer(cfg *config.Config) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Standard middlewares
	e.Use(middleware.TraceMiddleware())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"X-Trace-ID", "X-Span-ID"},
		AllowCredentials: true,
	}))

	healthHdl := health.NewHealthHandler()

	// AC6 - Health checks
	e.GET("/healthcheck", healthHdl.HealthCheck)
	e.GET("/healthcheck-readiness", healthHdl.Readiness)

	// API Documentation endpoint (/docs)
	e.GET("/docs", func(c echo.Context) error {
		html := `<!doctype html>
<html>
<head>
  <title>PortifyX API Reference</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <style>body { margin: 0; background-color: #0f172a; color: #f8fafc; font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; }</style>
</head>
<body>
  <div style="text-align: center;">
    <h1>PortifyX API Gateway</h1>
    <p>Edge Gateway routing traffic across PortifyX microservices.</p>
    <p><code>GET /healthcheck</code> | <code>GET /healthcheck-readiness</code> | <code>/api/v1/portfolios</code></p>
  </div>
</body>
</html>`
		return c.HTML(http.StatusOK, html)
	})

	// Upstream proxy routes
	setupProxyRoutes(e, cfg)

	return &Server{
		echo: e,
		cfg:  cfg,
	}
}

func setupProxyRoutes(e *echo.Echo, cfg *config.Config) {
	portfolioURL, err := url.Parse(cfg.Upstream.PortfolioServiceURL)
	if err == nil {
		balancer := echoMw.NewRoundRobinBalancer([]*echoMw.ProxyTarget{{URL: portfolioURL}})
		e.Group("/api/v1/portfolios", echoMw.Proxy(balancer))
	}

	exchangeURL, err := url.Parse(cfg.Upstream.ExchangeServiceURL)
	if err == nil {
		balancer := echoMw.NewRoundRobinBalancer([]*echoMw.ProxyTarget{{URL: exchangeURL}})
		e.Group("/api/v1/exchanges", echoMw.Proxy(balancer))
	}

	marketURL, err := url.Parse(cfg.Upstream.MarketDataServiceURL)
	if err == nil {
		balancer := echoMw.NewRoundRobinBalancer([]*echoMw.ProxyTarget{{URL: marketURL}})
		e.Group("/api/v1/market", echoMw.Proxy(balancer))
	}
}

func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}
