package api

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sam-ctd/job-tracker-agent/internal/store"
)

type Server struct {
	store       *store.Store
	pool        *pgxpool.Pool
	apiKey      string
	corsOrigins []string
}

func New(st *store.Store, pool *pgxpool.Pool, apiKey string, corsOrigins []string) *Server {
	return &Server{
		store:       st,
		pool:        pool,
		apiKey:      apiKey,
		corsOrigins: corsOrigins,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ready", s.handleReady)

	api := http.NewServeMux()
	api.HandleFunc("GET /applications", s.handleListApplications)
	api.HandleFunc("POST /applications", s.handleCreateApplication)
	api.HandleFunc("GET /applications/due", s.handleDue)
	api.HandleFunc("GET /applications/{id}", s.handleGetApplication)
	api.HandleFunc("PATCH /applications/{id}", s.handlePatchApplication)
	api.HandleFunc("DELETE /applications/{id}", s.handleDeleteApplication)
	api.HandleFunc("GET /stats", s.handleStats)

	mux.Handle("/api/", http.StripPrefix("/api", s.withAuth(api)))

	return s.withCORS(securityHeaders(mux))
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.pool.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
