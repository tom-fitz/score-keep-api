package imports

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Service struct {
	ctx     context.Context
	logger  *logrus.Logger
	version int
	db      *sql.DB
}

func NewService(ctx context.Context, logger *logrus.Logger, version int, db *sql.DB) *Service {
	h := &Service{
		ctx:     ctx,
		logger:  logger,
		version: version,
		db:      db,
	}
	return h
}
