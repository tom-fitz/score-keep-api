package league

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1/league/create", h.createLeagues)
	router.GET("/v1/league", h.getLeagues)
	router.GET("/v1/league/:id", h.getLeagueById)
}
