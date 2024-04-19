package calendar

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/calendar/v3"
)

type Handler struct {
	ctx     context.Context
	logger  *logrus.Logger
	version int
	db      *sql.DB
	gcp     *calendar.Service
}

func NewHandler(ctx context.Context, logger *logrus.Logger, version int, db *sql.DB, gcp *calendar.Service) *Handler {
	return &Handler{
		ctx:     ctx,
		logger:  logger,
		version: version,
		db:      db,
		gcp:     gcp,
	}
}
