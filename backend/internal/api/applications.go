package api

import (
	"encoding/json"
	"net/http"

	"github.com/sam-ctd/job-tracker-agent/internal/models"
)

func (s *Server) handleListApplications(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	company := r.URL.Query().Get("company")
	apps, err := s.store.List(r.Context(), status, company)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

func (s *Server) handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	var in models.ApplicationCreate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	app, err := s.store.Create(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, app)
}

func (s *Server) handleGetApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	app, err := s.store.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, app)
}

func (s *Server) handlePatchApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.ApplicationPatch
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	app, err := s.store.Patch(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, app)
}

func (s *Server) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDue(w http.ResponseWriter, r *http.Request) {
	apps, err := s.store.Due(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.Stats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
