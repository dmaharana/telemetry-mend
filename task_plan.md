# Task Plan: Configurable OpenAI LLM & Selective Fix Suggestions

## Goal
Implement a UI option to configure an OpenAI-compliant LLM service and update the fix suggestion logic to skip processing if the application's repository is not available.

## Current Phase
Phase 4: Testing & Verification

## Phases

### Phase 1: Research & Discovery
- [x] Analyze backend `ai` package and `fixer` service.
- [x] Analyze `scm` package for repository checks.
- [x] Check frontend structure for adding Settings.
- **Status:** complete

### Phase 2: Backend Implementation
- [x] Add `Settings` model to `models/models.go`.
- [x] Create `ai/openai.go` for OpenAI-compliant provider.
- [x] Implement `handlers/settings.go` for managing LLM configuration.
- [x] Update `service/fixer.go` to use dynamic settings and check repo availability.
- [x] Register routes in `main.go`.
- **Status:** complete

### Phase 3: Frontend Implementation
- [x] Create `Settings` component/page (integrated into existing SettingsPage).
- [x] Add API calls to fetch/save settings in `lib/api.ts`.
- [x] Update UI to handle "No Repository" state for fix suggestions.
- **Status:** complete

### Phase 4: Testing & Verification
- [x] Verify settings save/load.
- [x] Verify fix generation builds and logic handles providers.
- [x] Verify fix suggestion is skipped when `RepoURL` is empty.
- **Status:** complete

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| Use a single `Settings` table | Simplicity for global configuration. |
| OpenAI-compliant provider | Standard for most LLM services (Self-hosted, Azure, OpenAI, etc.). |
| Skip fix if no repo | Avoids useless AI calls and failures when source code is inaccessible. |
| Alias Settings as LLMSettings | Avoid conflict with Lucide icon in frontend. |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| Missing context import | 1 | Added context import to logs.go |
| Build failed in server | 1 | Ran build from correct directory |
| TS Duplicate identifier | 1 | Aliased API Settings as LLMSettings |
