package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pgreader/services"
)

type HTTPHandler struct {
	articles *services.ArticleService
}

func NewHTTPHandler(articles *services.ArticleService) *HTTPHandler {
	return &HTTPHandler{articles: articles}
}

func (h *HTTPHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /articles", h.listArticles)
	mux.HandleFunc("GET /articles/{id}", h.getArticle)
	mux.HandleFunc("PATCH /articles/{id}/read", h.setReadStatus)
}

func (h *HTTPHandler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HTTPHandler) listArticles(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.articles.List())
}

func (h *HTTPHandler) getArticle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}

	article, err := h.articles.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *HTTPHandler) setReadStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}

	var payload struct {
		IsRead bool `json:"isRead"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	article, err := h.articles.SetRead(id, payload.IsRead)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, article)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
