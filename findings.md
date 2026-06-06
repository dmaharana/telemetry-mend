# Findings - TelemetryMend

## Project Overview
TelemetryMend is an application designed to analyze application logs for failures, cluster them, and review source code from SCM to generate fix solutions using GenAI.

## Tech Stack
- **Backend:** Golang
- **ORM:** Bun (supporting SQLite for dev, PostgreSQL for prod)
- **Frontend:** React + Vite + TypeScript
- **UI Components:** shadcn/ui (Tailwind CSS)
- **Database:** SQLite (local development), PostgreSQL (production)

## Current Workspace State
- `telemetry-mend-server/`: Contains only `go.mod`.
- `telemetry-mend-ui/`: Boilerplate Vite + React app with shadcn/ui components already installed in `src/components/ui/`.

## Requirements Analysis
- Support multiple applications/microservices.
- Ingest logs with `app_id`, `environment`, `commit_hash`/`branch`, and `log_body`.
- Cluster logs by stripping dynamic variables.
- Fetch source code context from SCM (Git).
- Generate fixes using LLM.
- UI for triage and reviewing suggested fixes.
