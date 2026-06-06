# Progress: Configurable OpenAI LLM & Selective Fix Suggestions

## 2026-06-06
- Added `Settings` model and migration to backend.
- Implemented `OpenAIProvider` for standardized LLM access.
- Created `SettingsHandler` for managing LLM configuration via API.
- Refactored `FixerService` to:
    - Load LLM settings dynamically.
    - Skip fix suggestions if `RepoURL` is missing.
    - Use the chosen AI provider.
- Updated frontend `lib/api.ts` with settings types and functions.
- Updated `App.tsx` with:
    - LLM configuration UI in `SettingsPage`.
    - Repository availability check in `ClusterDetail`.
    - Toast notifications for AI actions.
- Verified both backend and frontend builds.
