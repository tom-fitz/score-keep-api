package league

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1/league/create", h.createLeagues)
	router.POST("/v1/league/:id/update", h.updateLeague)
	router.GET("/v1/league", h.getLeagues)
	router.GET("/v1/league/:id", h.getLeagueById)
	router.GET("/v1/league/:id/players", h.getPlayersByLeagueId)
	router.GET("/v1/league/:id/teams", h.getTeamsByLeagueId)
}
