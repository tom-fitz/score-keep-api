package imports

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (h *Service) RegisterRoutes(router *gin.Engine) {
	v := fmt.Sprintf("/v%d", h.version)
	router.POST(fmt.Sprintf("%s/league/:leagueId/teams/import", v), h.ImportTeams)
	router.POST(fmt.Sprintf("%s/league/:leagueId/players/import", v))
	router.GET("/import/healthcheck", h.healthcheck)
}

func (h *Service) healthcheck(c *gin.Context) {
	h.logger.Printf("healthcheck request received")
}
