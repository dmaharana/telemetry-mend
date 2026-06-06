package handlers

import (
	"encoding/json"
	"net/http"
	"telemetry-mend-server/internal/models"

	"github.com/uptrace/bun"
)

type SettingsHandler struct {
	db *bun.DB
}

func NewSettingsHandler(db *bun.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	settings := new(models.Settings)
	err := h.db.NewSelect().Model(settings).Limit(1).Scan(r.Context())
	if err != nil {
		// Return default settings if none exist
		settings = &models.Settings{
			LLMProvider: "mock",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var settings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing := new(models.Settings)
	err := h.db.NewSelect().Model(existing).Limit(1).Scan(r.Context())
	
	if err != nil {
		// Create new
		_, err = h.db.NewInsert().Model(&settings).Exec(r.Context())
	} else {
		// Update existing
		settings.ID = existing.ID
		_, err = h.db.NewUpdate().Model(&settings).WherePK().Exec(r.Context())
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
