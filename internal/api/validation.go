package api

import (
	"fmt"
	"net/http"
	"strings"
)

func RequirePathID(r *http.Request) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || strings.TrimSpace(parts[len(parts)-1]) == "" {
		return "", fmt.Errorf("path id is required")
	}
	return parts[len(parts)-1], nil
}

func ValidateBatchQuery(r *http.Request) error {
	batch := strings.TrimSpace(r.URL.Query().Get("batch"))
	if batch == "" {
		return fmt.Errorf("batch query is required")
	}
	if len(batch) > 64 {
		return fmt.Errorf("batch query is too long")
	}
	return nil
}

func IsJSON(r *http.Request) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type"))), "application/json")
}
