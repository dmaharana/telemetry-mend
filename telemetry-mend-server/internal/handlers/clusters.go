package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"telemetry-mend-server/internal/models"
	"telemetry-mend-server/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

type ClusterHandler struct {
	db           *bun.DB
	fixerService *service.FixerService
}

func NewClusterHandler(db *bun.DB, fixer *service.FixerService) *ClusterHandler {
	return &ClusterHandler{
		db:           db,
		fixerService: fixer,
	}
}

func (h *ClusterHandler) List(w http.ResponseWriter, r *http.Request) {
	appIDStr := r.URL.Query().Get("app_id")
	var clusters []models.ErrorCluster

	q := h.db.NewSelect().Model(&clusters)
	if appIDStr != "" {
		appID, _ := strconv.ParseInt(appIDStr, 10, 64)
		q = q.Where("app_id = ?", appID)
	}

	if err := q.Order("last_seen DESC").Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusters)
}

func (h *ClusterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	cluster := new(models.ErrorCluster)

	err := h.db.NewSelect().
		Model(cluster).
		Relation("Application").
		Where("ec.id = ?", id).
		Scan(r.Context())
	if err != nil {
		http.Error(w, "cluster not found", http.StatusNotFound)
		return
	}

	// Also fetch fixes
	var fixes []models.SuggestedFix
	h.db.NewSelect().Model(&fixes).Where("cluster_id = ?", id).Scan(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cluster": cluster,
		"fixes":   fixes,
	})
}

func (h *ClusterHandler) GenerateFix(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	fix, err := h.fixerService.GenerateAndStoreFix(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fix)
}
