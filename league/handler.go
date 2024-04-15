package league

import (
	"database/sql"
	"log"
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
