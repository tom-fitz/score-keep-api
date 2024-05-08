package imports

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	v := fmt.Sprintf("/v%d", h.version)
	router.POST(fmt.Sprintf("%s/league/:id/import", v), h.importLeagues)
	router.POST(fmt.Sprintf("%s/league/:leagueId/teams/import", v), h.ImportTeams)
	router.POST(fmt.Sprintf("%s/league/:leagueId/players/import", v))
	router.GET("/import/healthcheck", h.healthcheck)
}

func (h *Handler) healthcheck(c *gin.Context) {
	h.logger.Printf("healthcheck request received")
}
