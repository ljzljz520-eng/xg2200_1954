package api

import (
	"net/http"
	"strings"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = "local-request"
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func WithJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Header.Get("Content-Type") == "" {
			writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "content type required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Chain(handler http.Handler) http.Handler {
	return WithRequestID(WithJSON(handler))
}
