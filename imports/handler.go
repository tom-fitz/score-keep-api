package imports

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	ctx     context.Context
	logger  *logrus.Logger
	version int
	db      *sql.DB
}

func NewHandler(ctx context.Context, logger *logrus.Logger, version int, db *sql.DB) *Handler {
	h := &Handler{
		ctx:     ctx,
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
