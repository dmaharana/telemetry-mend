package handlers

import (
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
}

func (h *LogHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var req IngestLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.AppID == 0 || req.LogBody == "" {
		http.Error(w, "app_id and log_body are required", http.StatusBadRequest)
		return
	}

	// 1. Sanitize and Fingerprint
	sanitized := clustering.SanitizeLog(req.LogBody)
	fingerprint := clustering.GenerateFingerprint(sanitized)

	ctx := r.Context()

	// 2. Find or Create Cluster
	cluster := new(models.ErrorCluster)
	err := h.db.NewSelect().
		Model(cluster).
		Where("app_id = ? AND fingerprint = ?", req.AppID, fingerprint).
		Scan(ctx)

	if err != nil {
		// Create new cluster
		cluster = &models.ErrorCluster{
			AppID:       req.AppID,
			Fingerprint: fingerprint,
			LogTemplate: sanitized,
			Count:       1,
			LastSeen:    time.Now(),
		}
		_, err = h.db.NewInsert().Model(cluster).Exec(ctx)
		if err != nil {
			http.Error(w, "failed to create error cluster", http.StatusInternalServerError)
			return
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
			http.Error(w, "failed to update error cluster", http.StatusInternalServerError)
			return
		}
	}

	// 3. Store Log Entry
	logEntry := &models.LogEntry{
		AppID:       req.AppID,
		ClusterID:   cluster.ID,
		Environment: req.Environment,
		CommitHash:  req.CommitHash,
		LogBody:     req.LogBody,
	}
	_, err = h.db.NewInsert().Model(logEntry).Exec(ctx)
	if err != nil {
		http.Error(w, "failed to store log entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "accepted",
		"cluster_id": cluster.ID,
	})
}
