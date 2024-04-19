package league

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
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
