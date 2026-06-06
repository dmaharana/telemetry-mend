# Findings: Log Ingestion Research

## Vector Compatibility
- **HTTP Sink:** Vector can send logs to an HTTP endpoint using either `json` or `ndjson` (newline-delimited JSON) format.
- **Batching:** Vector typically batches logs and sends them as a JSON array `[...]`.
- **Headers:** Vector can be configured to send custom headers, which is perfect for authentication (e.g., `X-Ingest-Token`).

## Standardized Payload
To be compatible with most loggers and shippers, we should accept a payload like:
```json
{
  "timestamp": "2023-10-27T10:00:00Z",
  "message": "Error connecting to database",
  "level": "error",
  "service": "my-app",
  "metadata": {
    "commit_hash": "abc1234",
    "environment": "production"
  }
}
```
Or a simplified version that matches our current model but allows for extra fields:
```json
{
  "log_body": "...",
  "commit_hash": "...",
  "environment": "...",
  "app_id": 123
}
```

## Authentication Strategy
- Each `Application` in the database will have a unique `api_key` or `ingest_token`.
- The token will be passed in the `Authorization: Bearer <token>` or `X-Ingest-Token` header.
- This allows the server to identify the `app_id` without requiring it in the JSON body, which is safer and more flexible.
