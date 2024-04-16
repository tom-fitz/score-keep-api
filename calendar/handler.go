package calendar

import (
	"context"
	"database/sql"
	"google.golang.org/api/calendar/v3"
	"log"
)

type Handler struct {
	ctx     context.Context
	logger  *log.Logger
	version int
	db      *sql.DB
	gcp     *calendar.Service
}

func NewHandler(ctx context.Context, logger *log.Logger, version int, db *sql.DB, gcp *calendar.Service) *Handler {
	return &Handler{
		ctx:     ctx,
		logger:  logger,
		version: version,
		db:      db,
		gcp:     gcp,
	}
}
