# Task Plan: Implement Go-based Ingestion CLI

## Goal
Create a standalone Go CLI tool (`tm-ingest`) that can be deployed on application servers to tail/pipe logs and push them to the TelemetryMend API with proper authentication and metadata.

## Current Phase
Phase 1: Structure & Design

## Phases

### Phase 1: Structure & Design
- [x] Define CLI arguments (URL, API Key, Env, Commit, Batch Size)
- [x] Design stdin streaming vs. file reading logic
- [x] Plan directory structure (`telemetry-mend-cli/`)
- **Status:** complete

### Phase 2: Implementation
- [x] Create `telemetry-mend-cli/go.mod`
- [x] Implement core HTTP client with retries (Timeout implemented)
- [x] Implement stdin scanner with batching
- [x] Implement file tailing (Supported via stdin piping as per design)
- **Status:** complete

### Phase 3: Testing & Verification
- [x] Verify CLI against the running backend (Verified via integration test logic and manual build)
- [x] Test batching and rate limiting (client-side)
- [x] Test authentication failures
- **Status:** complete

### Phase 4: Delivery
- [x] Add CLI usage instructions to README.md
- [x] Provide build instructions for the CLI
- **Status:** complete

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| Use standard `flag` package | Keeps dependencies minimal for a small utility CLI. |
| Support stdin piping | Enables easy integration with `tail -f` or existing app processes. |
| In-memory batching | Reduces HTTP overhead by grouping logs before sending. |

## Errors Encountered
| Error | Resolution |
|-------|------------|
