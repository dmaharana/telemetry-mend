package handlers

import (
	"encoding/json"
	"net/http"
	"telemetry-mend-server/internal/models"

	"github.com/uptrace/bun"
)

type AppHandler struct {
	db *bun.DB
}

func NewAppHandler(db *bun.DB) *AppHandler {
	return &AppHandler{db: db}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var app models.Application
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
