# portifyx-api-gateway

Edge API Gateway and reverse proxy service for the PortifyX platform, built on Go and Echo v4.

## Key Features

- **Edge Entry Point**: Unified entry point for client applications (portifyx-frontend, mobile, third-party clients).
- **CORS Handling**: Permissive browser CORS headers with trace header exposure (`X-Trace-ID`, `X-Span-ID`).
- **Trace Propagation**: Auto-injects and propagates distributed tracing IDs across HTTP upstreams.
- **Health Checks**:
  - `GET /healthcheck`: Liveness probe (returns `{"status":"ok"}`).
  - `GET /healthcheck-readiness`: Readiness probe.
- **Interactive Documentation**: Interactive API reference available at `/docs`.

## Running Locally

```bash
# Build and run with Cobra CLI
make run
# Or directly:
go run main.go api
```
