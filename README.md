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
