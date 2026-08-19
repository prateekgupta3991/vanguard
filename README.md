# Vanguard

Vanguard is the central control plane for AI agents. It will provide agent
registration, maker-checker permissions, and audit APIs for enforcement points
such as Sentinel.

The HTTP API is built with [Gin](https://gin-gonic.com/).

## Development

Vanguard requires Go 1.24 or newer.

```sh
go test ./...
go run ./cmd/vanguard
```

Run the admin UI in a second terminal:

```sh
cd web
npm install
npm run dev
```

The UI is available at `http://localhost:5173` and proxies API requests to
Vanguard at `http://localhost:8080`.

The server listens on `:8080` by default and exposes:

- `GET /healthz` — process health
- `GET /readyz` — service readiness
- `POST /api/v1/agents` — register an agent or refresh an existing registration
- `GET /api/v1/agents` — list registered agents
- `GET /api/v1/agents/:agentId` — retrieve a registered agent

Configuration is provided through environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `VANGUARD_HTTP_ADDRESS` | `:8080` | HTTP listen address |
| `VANGUARD_SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown timeout |
