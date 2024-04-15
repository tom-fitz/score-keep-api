package imports

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	logger  *log.Logger
	version int
	db      *sql.DB
}

func NewHandler(logger *log.Logger, version int, db *sql.DB) *Handler {
	return &Handler{logger, version, db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		h.importStatus(w, r)
	case http.MethodPost:
		h.importLeagues(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) importStatus(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ok",
		"method": "get",
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}
