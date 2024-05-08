package league

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type League struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

func (h *Handler) createLeagues(c *gin.Context) {
	var requestBody []map[string]interface{}

	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	var ids []int
	query := "INSERT INTO score_keep_db.public.leagues (name, level) VALUES ($1, $2) RETURNING id"
	for _, league := range requestBody {
		var id int
		err := h.db.QueryRow(query, league["name"], league["level"]).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert data into database: %v", err)})
			return
		}
		ids = append(ids, id)
	}

	c.JSON(http.StatusOK, gin.H{"ids": ids})
}

func (h *Handler) getLeagues(c *gin.Context) {
	rows, err := h.db.Query("SELECT * FROM score_keep_db.public.leagues")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error fetching leagues": fmt.Sprintf("%v", err)})
		return
	}
	defer rows.Close()

	var leagues []interface{}
	for rows.Next() {
		var league League
		if err := rows.Scan(&league.Id, &league.Name, &league.Level); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error scanning league": fmt.Sprintf("%v", err)})
			return
		}
		leagues = append(leagues, league)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error iterating over leagues": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leagues": leagues})
}

func (h *Handler) getLeagueById(c *gin.Context) {
	lid := c.Param("id")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	query := "SELECT * FROM score_keep_db.public.leagues WHERE id = $1"
	var league League
	if err := h.db.QueryRow(query, lid).Scan(&league.Id, &league.Name, &league.Level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error scanning league": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"league": league})
}
