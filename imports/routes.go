package imports

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1/league/:id/import", h.importLeagues)
	router.GET("/import/healthcheck", h.healthcheck)
}

func (h *Handler) healthcheck(c *gin.Context) {
	h.logger.Printf("healthcheck request received")
}
