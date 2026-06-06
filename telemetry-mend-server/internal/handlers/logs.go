package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"telemetry-mend-server/internal/clustering"
	"telemetry-mend-server/internal/models"
	"time"

	"github.com/uptrace/bun"
)

type LogHandler struct {
	db *bun.DB
}

func NewLogHandler(db *bun.DB) *LogHandler {
	return &LogHandler{db: db}
}

type IngestLogRequest struct {
	AppID       int64  `json:"app_id"`
	Environment string `json:"environment"`
	CommitHash  string `json:"commit_hash"`
	LogBody     string `json:"log_body"`
	Message     string `json:"message"` // Vector often uses 'message'
}

func (h *LogHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Authenticate
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.Header.Get("Authorization")
		if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
			apiKey = apiKey[7:]
		}
	}

	if apiKey == "" {
		http.Error(w, "missing API key", http.StatusUnauthorized)
		return
	}

	app := new(models.Application)
	err := h.db.NewSelect().
		Model(app).
		Where("api_key = ?", apiKey).
		Scan(ctx)
	if err != nil {
		http.Error(w, "invalid API key", http.StatusUnauthorized)
		return
	}

	// 2. Decode Body (detect single or array)
	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var logs []IngestLogRequest
	if raw[0] == '[' {
		if err := json.Unmarshal(raw, &logs); err != nil {
			http.Error(w, "failed to decode log array", http.StatusBadRequest)
			return
		}
	} else {
		var single IngestLogRequest
		if err := json.Unmarshal(raw, &single); err != nil {
			http.Error(w, "failed to decode log object", http.StatusBadRequest)
			return
		}
		logs = append(logs, single)
	}

	// 3. Process Logs
	var processedCount int
	for _, req := range logs {
		// Normalize log body
		logBody := req.LogBody
		if logBody == "" {
			logBody = req.Message
		}
		if logBody == "" {
			continue
		}

		// Use AppID from authenticated app if not provided or different
		appID := app.ID

		if err := h.processLog(ctx, appID, req.Environment, req.CommitHash, logBody); err != nil {
			// Log error but continue with others
			continue
		}
		processedCount++
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "accepted",
		"processed": processedCount,
	})
}

func (h *LogHandler) processLog(ctx context.Context, appID int64, env, commit, body string) error {
	// 1. Sanitize and Fingerprint
	sanitized := clustering.SanitizeLog(body)
	fingerprint := clustering.GenerateFingerprint(sanitized)

	// 2. Find or Create Cluster
	cluster := new(models.ErrorCluster)
	err := h.db.NewSelect().
		Model(cluster).
		Where("app_id = ? AND fingerprint = ?", appID, fingerprint).
		Scan(ctx)

	if err != nil {
		// Create new cluster
		cluster = &models.ErrorCluster{
			AppID:       appID,
			Fingerprint: fingerprint,
			LogTemplate: sanitized,
			Count:       1,
			LastSeen:    time.Now(),
		}
		_, err = h.db.NewInsert().Model(cluster).Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		// Update existing cluster
		cluster.Count++
		cluster.LastSeen = time.Now()
		_, err = h.db.NewUpdate().
			Model(cluster).
			Column("count", "last_seen").
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// 3. Store Log Entry
	logEntry := &models.LogEntry{
		AppID:       appID,
		ClusterID:   cluster.ID,
		Environment: env,
		CommitHash:  commit,
		LogBody:     body,
	}
	_, err = h.db.NewInsert().Model(logEntry).Exec(ctx)
	return err
}
