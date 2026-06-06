# Task Plan - TelemetryMend Implementation

## Goal
Implement the TelemetryMend application according to `spec.md`, enabling log-based error clustering and AI-powered source code remediation.

## Phases

### Phase 1: Backend Infrastructure & Models `[complete]`
- [x] Initialize Go project structure in `telemetry-mend-server`.
- [x] Set up Bun ORM with SQLite.
- [x] Define models: `Application`, `Deployment`, `ErrorCluster`, `LogEntry`, `SuggestedFix`.
- [x] Implement database migrations/initialization.

### Phase 2: Log Ingestion & Clustering `[complete]`
- [x] Create API endpoint for log ingestion.
- [x] Implement log sanitization (regex-based).
- [x] Implement clustering logic (fingerprinting).
- [x] Store clustered errors in database.

### Phase 3: SCM Integration `[complete]`
- [x] Implement Git client to interact with repositories.
- [x] Logic to fetch specific file context based on stack traces.

### Phase 4: AI Fix Generation `[complete]`
- [x] Integrate with LLM provider (mocked).
- [x] Design prompt for fix generation.
- [x] Implement structured output parsing for fixes.

### Phase 5: Frontend Development `[complete]`
- [x] Implement Application Management UI.
- [x] Implement Error Cluster Dashboard.
- [x] Implement Fix Review / Diff View UI.

### Phase 6: End-to-End Integration & Testing `[complete]`
- [x] Connect Frontend to Backend (via Vite proxy).
- [x] Verify full flow from log ingestion to fix generation (via integration test).

### Phase 7: Single Binary Packaging `[complete]`
- [x] Build frontend assets.
- [x] Embed assets using Go `embed` package.
- [x] Implement static file serving with SPA routing support.

## Progress Log
- [2026-06-06] Initialized project and planning files.
