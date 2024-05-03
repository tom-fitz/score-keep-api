package league

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
