package api

import (
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ReadyResponse struct {
	Service string `json:"service"`
	Ready   bool   `json:"ready"`
	Build   string `json:"build"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{Code: code, Message: message})
}

func WriteReady(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, ReadyResponse{Service: "drone-telemetry", Ready: true, Build: time.Unix(0, 0).UTC().Format("20060102")})
}

func MethodAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", "GET, POST")
	WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method is not supported")
}
