package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"telemetry-mend-server/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

type AppHandler struct {
	db *bun.DB
}

func NewAppHandler(db *bun.DB) *AppHandler {
	return &AppHandler{db: db}
}

func generateAPIKey() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var app models.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if app.APIKey == "" {
		app.APIKey = "tm_" + generateAPIKey()
	}

	if _, err := h.db.NewInsert().Model(&app).Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	var apps []models.Application
	if err := h.db.NewSelect().Model(&apps).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var app models.Application
	if err := h.db.NewSelect().Model(&app).Where("id = ?", id).Scan(r.Context()); err != nil {
		http.Error(w, "application not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *AppHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var existing models.Application
	if err := h.db.NewSelect().Model(&existing).Where("id = ?", id).Scan(r.Context()); err != nil {
		http.Error(w, "application not found", http.StatusNotFound)
		return
	}

	var app models.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Only update editable fields
	if app.Name != "" {
		existing.Name = app.Name
	}
	if app.RepoURL != "" {
		existing.RepoURL = app.RepoURL
	}
	if app.Language != "" {
		existing.Language = app.Language
	}
	if app.DefaultBranch != "" {
		existing.DefaultBranch = app.DefaultBranch
	}

	if _, err := h.db.NewUpdate().Model(&existing).WherePK().Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete related SuggestedFixes via ErrorClusters
	var clusterIDs []int64
	if err := tx.NewSelect().
		Table("error_clusters").
		Column("id").
		Where("app_id = ?", id).
		Scan(r.Context(), &clusterIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(clusterIDs) > 0 {
		if _, err := tx.NewDelete().
			Table("suggested_fixes").
			Where("cluster_id IN (?)", bun.In(clusterIDs)).
			Exec(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Delete related ErrorClusters
	if _, err := tx.NewDelete().
		Model((*models.ErrorCluster)(nil)).
		Where("app_id = ?", id).
		Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete related LogEntries
	if _, err := tx.NewDelete().
		Model((*models.LogEntry)(nil)).
		Where("app_id = ?", id).
		Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete related Deployments
	if _, err := tx.NewDelete().
		Model((*models.Deployment)(nil)).
		Where("app_id = ?", id).
		Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the Application itself
	if _, err := tx.NewDelete().
		Model((*models.Application)(nil)).
		Where("id = ?", id).
		Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
