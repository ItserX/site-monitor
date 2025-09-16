package crud

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"site-monitor/internal/storage"
	"site-monitor/pkg/logger"
)

type Handler struct {
	storage storage.Storage
	log     *logger.Logger
}

func NewHandler(storage storage.Storage, log *logger.Logger) *Handler {
	return &Handler{storage: storage, log: log}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/sites", h.handleGetSites)
	r.Get("/sites/{id}", h.handleGetSiteByID)
	r.Post("/sites", h.handleAddSite)
	r.Put("/sites/{id}", h.handleUpdateSite)
	r.Delete("/sites/{id}", h.handleDeleteSite)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func (h *Handler) handleGetSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.storage.GetSites(r.Context())
	if err != nil {
		h.log.Sugar.Errorw("Failed to get sites", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Sugar.Infow("Fetched all sites", "count", len(sites))
	writeJSON(h.log, w, sites, http.StatusOK)
}

func (h *Handler) handleGetSiteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	site, err := h.storage.GetSiteByID(r.Context(), id)
	if err != nil {
		h.log.Sugar.Errorw("Failed to get site by ID", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if site == nil {
		h.log.Sugar.Warnw("Site not found", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.log.Sugar.Infow("Fetched site by ID", "id", id)
	writeJSON(h.log, w, site, http.StatusOK)
}

func (h *Handler) handleAddSite(w http.ResponseWriter, r *http.Request) {
	var site storage.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		h.log.Sugar.Warnw("Invalid request body for AddSite", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.storage.AddSite(r.Context(), site)
	if err != nil {
		h.log.Sugar.Errorw("Failed to add site", "url", site.URL, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	site.ID = id
	h.log.Sugar.Infow("Site added", "id", id, "url", site.URL)
	writeJSON(h.log, w, site, http.StatusCreated)
}

func (h *Handler) handleUpdateSite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var site storage.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		h.log.Sugar.Warnw("Invalid request body for UpdateSite", "id", id, "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	site.ID = id

	if err := h.storage.UpdateSite(r.Context(), site); err != nil {
		h.log.Sugar.Errorw("Failed to update site", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Sugar.Infow("Site updated", "id", id, "url", site.URL, "active", site.Active)
	writeJSON(h.log, w, site, http.StatusOK)
}

func (h *Handler) handleDeleteSite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.storage.DeleteSite(r.Context(), id); err != nil {
		h.log.Sugar.Errorw("Failed to delete site", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Sugar.Infow("Site deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(log *logger.Logger, w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Sugar.Errorw("Failed to write JSON response", "error", err)
	}
}
