package imports

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
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

func (h *Handler) importStatus(c *gin.Context) {
	resp := map[string]string{
		"status": "ok",
		"method": "get",
	}
	c.JSON(http.StatusOK, resp)
}
