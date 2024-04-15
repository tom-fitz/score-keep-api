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
	h := &Handler{
		logger:  logger,
		version: version,
		db:      db,
	}
	return h
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
