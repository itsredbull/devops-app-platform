package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/itsredbull/devops-app-platform/app/internal/metrics"
	"github.com/itsredbull/devops-app-platform/app/internal/store"
)

type Server struct {
	store *store.PostgresStore
}

func NewServer(st *store.PostgresStore) *Server {
	return &Server{store: st}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /metrics", MetricsHandler())

	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /readyz", s.readyz)

	mux.HandleFunc("POST /api/v1/targets", s.createTarget)
	mux.HandleFunc("GET /api/v1/targets", s.listTargets)
	mux.HandleFunc("GET /api/v1/targets/{id}", s.getTarget)
	mux.HandleFunc("DELETE /api/v1/targets/{id}", s.deleteTarget)
	mux.HandleFunc("GET /api/v1/checks", s.listChecks)
	mux.HandleFunc("GET /api/v1/status", s.status)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "time": time.Now().UTC()})
}

func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.store.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database not ready")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ready", "time": time.Now().UTC()})
}

func (s *Server) createTarget(w http.ResponseWriter, r *http.Request) {
	var in store.CreateTargetInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	in.URL = strings.TrimSpace(in.URL)
	if in.URL == "" || !(strings.HasPrefix(in.URL, "http://") || strings.HasPrefix(in.URL, "https://")) {
		writeError(w, http.StatusBadRequest, "url must start with http:// or https://")
		return
	}
	if in.CheckIntervalSeconds <= 0 {
		in.CheckIntervalSeconds = 30
	}
	if in.TimeoutSeconds <= 0 {
		in.TimeoutSeconds = 10
	}

	target, err := s.store.CreateTarget(r.Context(), in)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			writeError(w, http.StatusConflict, "target URL already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create target")
		return
	}

	writeJSON(w, http.StatusCreated, target)
}

func (s *Server) listTargets(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListTargets(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list targets")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) getTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.store.GetTarget(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "target not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get target")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) deleteTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	deleted, err := s.store.DeleteTarget(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete target")
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, "target not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listChecks(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("target_id")
	limit := 100
	if qs := r.URL.Query().Get("limit"); qs != "" {
		if n, err := strconv.Atoi(qs); err == nil {
			limit = n
		}
	}

	items, err := s.store.ListChecks(r.Context(), targetID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list checks")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListLatestStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch status")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func WithMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		metrics.RecordAPIRequest(r.Method, r.URL.Path, rec.status)
	})
}

func MetricsHandler() http.Handler {
	return metrics.Handler()
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"error": msg})
}
