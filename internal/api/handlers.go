package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"telemetry.local/drone/internal/domain"
	"telemetry.local/drone/internal/service"
)

type Server struct {
	Service *service.Service
	Handler http.Handler
}

func NewServer(s *service.Service) *Server {
	server := &Server{Service: s}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.health)
	mux.HandleFunc("/records", server.records)
	mux.HandleFunc("/import", server.importRecords)
	mux.HandleFunc("/review", server.review)
	server.Handler = mux
	return server
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "drone-telemetry"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	status := domain.Status(r.URL.Query().Get("status"))
	records, err := s.Service.Search(r.URL.Query().Get("q"), status)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) importRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	batch := r.URL.Query().Get("batch")
	if strings.TrimSpace(batch) == "" {
		batch = "http"
	}
	result, summary, err := s.Service.ImportBatch(strings.NewReader(string(data)), batch)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"accepted": len(result.Records), "rejected": result.Rejected, "summary": summary})
}

func (s *Server) review(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var request struct {
		ID     string `json:"id"`
		Action string `json:"action"`
		Title  string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var (
		record domain.Record
		err    error
	)
	switch strings.ToLower(request.Action) {
	case "archive":
		record, err = s.Service.ReviewAndArchive(request.ID)
	case "publish":
		record, err = s.Service.PublishTitle(request.ID, request.Title)
	default:
		err = &badActionError{action: request.Action}
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

type badActionError struct{ action string }

func (e *badActionError) Error() string { return "unsupported review action: " + e.action }

func Handler(s *service.Service) http.Handler {
	return NewServer(s).Handler
}
