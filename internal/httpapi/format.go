package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"warrantyservice/internal/model"
)

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

func QueryParams(r *http.Request) (string, string, string) {
	query := r.URL.Query()
	format := strings.ToLower(strings.TrimSpace(query.Get("format")))
	if format == "" {
		format = "html"
	}
	return strings.TrimSpace(query.Get("phone")), strings.TrimSpace(query.Get("serial")), format
}

func EncodeResult(result model.QueryResult) ([]byte, error) {
	return json.Marshal(struct {
		Data   model.QueryResult `json:"data"`
		Status string            `json:"status"`
	}{Data: result, Status: result.Status.String()})
}

func WriteError(w http.ResponseWriter, status int, code, message, hint string) {
	writeJSON(w, status, ErrorPayload{Code: code, Message: message, Hint: hint})
}

func ParseBool(value string, fallback bool) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	if value == "1" || value == "true" || value == "yes" {
		return true
	}
	return false
}

func ParseID(value string) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("id is required")
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("id must be positive")
	}
	return id, nil
}

func AcceptsJSON(r *http.Request) bool {
	value := strings.ToLower(r.Header.Get("Accept"))
	return strings.Contains(value, "application/json") || r.URL.Query().Get("format") == "json"
}
func SetNoCache(w http.ResponseWriter) { w.Header().Set("Cache-Control", "no-store") }
func SetRequestID(w http.ResponseWriter, value string) {
	if value != "" {
		w.Header().Set("X-Request-ID", value)
	}
}

func MethodAllowed(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	return false
}

func DecodeJSONBody(r *http.Request, target any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
