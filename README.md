# TelemetryMend

TelemetryMend is an automated error analysis and remediation tool. It clusters application logs into unique failure patterns and uses GenAI to suggest source code fixes by analyzing stack traces and fetching relevant context from SCM.

## 🚀 Features

-   **Log Ingestion & Clustering:** Automatically sanitizes and fingerprints logs to identify unique error patterns.
-   **SCM Integration:** Parses stack traces to fetch source code context directly from Git repositories.
-   **AI-Powered Remediation:** Generates targeted fixes and diffs using LLM analysis (Mocked in current version).
-   **Unified Dashboard:** A modern React interface for triaging clusters and reviewing suggested fixes.
-   **Single Binary:** The entire application (Frontend + Backend) is packaged into a single executable.

## 🛠️ Tech Stack

-   **Backend:** Go 1.26+, Bun ORM, Chi Router.
-   **Frontend:** React, TypeScript, Vite, shadcn/ui, TanStack Query.
-   **Database:** SQLite (default for development).

## 📦 Building

To build the single binary application, you need `go` and `pnpm` installed.

```bash
# Build the entire application
make build
```

The resulting binary will be located at `bin/telemetry-mend`.

## 🏃 Running

```bash
# Run the application
./bin/telemetry-mend
```

The server will start at `http://localhost:8080`.

## 🛠️ Development

### Prerequisites
- Go 1.26+
- Node.js & pnpm

### Backend
```bash
cd telemetry-mend-server
go run ./cmd/server/main.go
```

### Frontend
```bash
cd telemetry-mend-ui
pnpm install
pnpm dev
```
(Frontend dev server proxies `/api` to `http://localhost:8080` by default).

## 🧪 Testing

```bash
cd telemetry-mend-server
go test ./...
```

## 📡 Log Ingestion API

TelemetryMend provides a robust API for ingesting logs from your applications or log shippers like Vector.

### Authentication
Every application is assigned a unique API key upon creation. Include this key in the header of your requests:
- `X-API-Key: tm_your_api_key`
- OR `Authorization: Bearer tm_your_api_key`

### Endpoint: `POST /api/logs/ingest`

The API supports both single log entries and batched arrays.

#### Single Log Entry
```bash
curl -X POST http://localhost:8080/api/logs/ingest \
  -H "X-API-Key: tm_your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "environment": "production",
    "commit_hash": "a1b2c3d4",
    "message": "panic: runtime error: index out of range"
  }'
```

#### Batched Logs (JSON Array)
```bash
curl -X POST http://localhost:8080/api/logs/ingest \
  -H "X-API-Key: tm_your_api_key" \
  -H "Content-Type: application/json" \
  -d '[
    {"message": "Error A", "environment": "prod"},
    {"message": "Error B", "environment": "staging"}
  ]'
```

### 📦 Integration with Vector
Example `vector.toml` configuration to ship logs to TelemetryMend:

```toml
[sinks.telemetry_mend]
type = "http"
inputs = ["your_log_source"]
uri = "http://localhost:8080/api/logs/ingest"
method = "post"
encoding.codec = "json"

[sinks.telemetry_mend.headers]
X-API-Key = "tm_your_api_key"
```

## 🛠️ Ingestion CLI (`tm-ingest`)

TelemetryMend includes a lightweight Go-based CLI tool for pushing logs from application servers.

### Installation
Build the CLI from the project root:
```bash
make build-cli
```
The binary will be available at `bin/tm-ingest`.

### Usage
The CLI reads logs from **stdin**, allowing you to pipe output from any process or log file.

#### Piping from an application
```bash
./my-app | ./bin/tm-ingest -key tm_your_api_key -env production
```

#### Tailing a log file
```bash
tail -f /var/log/app.log | ./bin/tm-ingest -key tm_your_api_key -env production
```

#### Configuration
| Flag | Env Var | Default | Description |
| :--- | :--- | :--- | :--- |
| `-key` | `TM_API_KEY` | - | **Required.** Your application API Key. |
| `-url` | `TM_URL` | `http://localhost:8080/api/logs/ingest` | Backend ingestion URL. |
| `-env` | - | `production` | Environment tag (e.g. prod, staging). |
| `-commit` | - | - | Optional git commit hash. |
| `-batch` | - | `10` | Number of logs to group before sending. |
| `-interval`| - | `5s` | Max time to wait before flushing a partial batch. |

