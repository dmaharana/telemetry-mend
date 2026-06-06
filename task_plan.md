# Task Plan: Refine Log Ingestion API for Simplicity and Shipper Compatibility

## Goal
Implement a robust REST API for log ingestion that supports both direct application pushes and batch shipping from agents like Vector, including authentication and standardized metadata handling.

## Current Phase
Phase 5: Delivery

## Phases

### Phase 1: Requirements & Discovery
- [x] Identify major log ingestion methods (Direct, Agent, Queue, SDK)
- [x] Get user preference for recommended approach (REST API with Shipper compatibility)
- [x] Research Vector HTTP sink compatibility
- **Status:** complete

### Phase 2: Planning & Structure
- [x] Define standardized log ingestion payload (Single vs Batch)
- [x] Design authentication mechanism (Ingest Token)
- [x] Plan model updates for `Application` to support tokens
- **Status:** complete

### Phase 3: Implementation
- [x] Update `models.Application` to include `api_key` or `ingest_token`
- [x] Update `handlers.LogHandler` to handle batching and authentication
- [x] Implement metadata extraction from standardized fields
- **Status:** complete

### Phase 4: Testing & Verification
- [x] Create a test script to simulate direct single-log ingestion (in integration_test.go)
- [x] Create a test script to simulate Vector-style batch ingestion (in integration_test.go)
- [x] Verify logs are correctly clustered and mapped to applications
- **Status:** complete

### Phase 5: Delivery
- [x] Document the new ingestion format and authentication
- [x] Provide example Vector configuration
- [x] Add examples and usage patterns to README.md
- **Status:** complete

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| Start with REST API | Balances simplicity with immediate functionality as requested by user. |
| Shipper Compatibility | Ensures the system can scale to production agents like Vector/Fluent Bit. |
| Token-based Auth | Essential for security in a multi-application environment. |
| Automatic API Key Gen | Improves DX by making apps ready for ingestion immediately. |

## Errors Encountered
| Error | Resolution |
|-------|------------|
| undefined: context in logs.go | Added "context" to imports in logs.go. |
