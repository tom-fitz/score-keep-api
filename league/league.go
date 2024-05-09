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

type TeamInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Captain    string `json:"captain"`
	FirstYear  int    `json:"firstyear"`
	LeagueID   int    `json:"league_id"`
	LeagueName string `json:"league_name"`
	TeamID     int    `json:"team_id"`
	TeamName   string `json:"team_name"`
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

func (h *Handler) updateLeague(c *gin.Context) {
	lid := c.Param("id")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}
	var requestBody League
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	league, err := validateLeague(lid, h.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
		return
	}
	league.Name = requestBody.Name
	league.Level = requestBody.Level

	query := `UPDATE score_keep_db.public.leagues SET name = $1, level = $2 WHERE id = $3`

	if err = h.db.QueryRow(query, league.Name, league.Level, league.Id).Scan(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating league in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "League updated successfully"})
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

func (h *Handler) getTeamsByLeagueId(c *gin.Context) {
	lid := c.Param("id")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	_, err := validateLeague(lid, h.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
		return
	}

	query := `SELECT t.id, t.name, t.captain, t.firstyear, lt.league_id, lt.league_name, lt.team_id, lt.team_name 
              FROM score_keep_db.public.teams t 
              JOIN score_keep_db.public.league_team lt ON t.id = lt.team_id
              WHERE lt.league_id = $1`

	rows, err := h.db.Query(query, lid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error querying teams": fmt.Sprintf("%v", err)})
		return
	}
	defer rows.Close()

	var teams []TeamInfo
	for rows.Next() {
		var team TeamInfo
		if err := rows.Scan(&team.ID, &team.Name, &team.Captain, &team.FirstYear, &team.LeagueID, &team.LeagueName, &team.TeamID, &team.TeamName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error scanning team": fmt.Sprintf("%v", err)})
			return
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error iterating rows": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

type PlayerInfo struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Level     int    `json:"level"`
	Phone     string `json:"phone"`
	Usanum    string `json:"usanum"`
	TeamName  string `json:"teamName"`
	TeamID    int    `json:"teamId"`
}

func (h *Handler) getPlayersByLeagueId(c *gin.Context) {
	lid := c.Param("id")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	_, err := validateLeague(lid, h.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
		return
	}

	query := `SELECT p.id, p.email, p.firstname, p.lastname, p.level, p.phone, p.usanum, pt.team_name, pt.team_id
              FROM score_keep_db.public.players p
              JOIN score_keep_db.public.player_team pt ON p.usanum = pt.usanum
              JOIN score_keep_db.public.league_team lt ON pt.team_id = lt.team_id
              WHERE lt.league_id = $1`

	rows, err := h.db.Query(query, lid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error querying players": fmt.Sprintf("%v", err)})
		return
	}
	defer rows.Close()

	var players []PlayerInfo
	for rows.Next() {
		var player PlayerInfo
		if err := rows.Scan(&player.ID, &player.Email, &player.FirstName, &player.LastName, &player.Level, &player.Phone, &player.Usanum, &player.TeamName, &player.TeamID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error scanning player": fmt.Sprintf("%v", err)})
			return
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error iterating rows": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})
}
