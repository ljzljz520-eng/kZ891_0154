package httpapi

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"warrantyservice/internal/config"
	"warrantyservice/internal/model"
	"warrantyservice/internal/warranty"
)

type Handler struct {
	service  *warranty.Service
	config   config.Config
	template *template.Template
}

func NewHandler(service *warranty.Service, cfg config.Config) *Handler {
	tmpl := template.Must(template.New("query").Parse(queryPage))
	return &Handler{service: service, config: cfg, template: tmpl}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "" {
		h.Home(w, r)
		return
	}
	if path == "/healthz" {
		h.Health(w, r)
		return
	}
	if path == "/warranty/query" {
		h.Query(w, r)
		return
	}
	if path == "/admin/warranties" && r.Method == http.MethodPost {
		h.AdminUpsert(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.template.Execute(w, map[string]any{"ServiceName": h.config.ServiceName, "Hotline": h.config.SupportHotline})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": h.config.ServiceName})
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	phone, serial := r.URL.Query().Get("phone"), r.URL.Query().Get("serial")
	result, err := h.service.Query(r.Context(), phone, serial)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if wantsJSON(r) {
		writeJSON(w, result.Status.HTTPCode(), result)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.template.Execute(w, map[string]any{"ServiceName": h.config.ServiceName, "Hotline": h.config.SupportHotline, "Phone": phone, "Serial": serial, "Result": result, "Found": result.Record != nil, "Valid": result.Status.IsSuccess()})
}

func (h *Handler) AdminUpsert(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+h.config.AdminToken {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "管理员令牌无效"})
		return
	}
	var input warranty.UpsertInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式无效"})
		return
	}
	operator := r.Header.Get("X-Operator")
	if operator == "" {
		operator = "admin"
	}
	record, err := h.service.Upsert(r.Context(), input, operator)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "saved", "record": record})
}

func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json") || r.URL.Query().Get("format") == "json"
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

var _ model.QueryResult
