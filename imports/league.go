package imports

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (h *Handler) importLeagues(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	if err := validateLeague(id, h.db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
		return
	}

	if c.Request.MultipartForm == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format. Expected multipart/form-data"})
		return
	}

	// Set max memory for file uploads (10 MB)
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ParseMultipartForm: %v", err.Error())})
		return
	}

	//if c.Request.Header.Get("Content-Type") != "multipart/form-data" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %s. Expected multipart/form-data", c.Request.Header.Get("Content-Type"))})
	//	return
	//}
	//
	//// Set max memory for file uploads (10 MB)
	//err := c.Request.ParseMultipartForm(10 << 20)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ParseMultipartForm: %v", err.Error())})
	//	return
	//}

	teamsFile, _, err := c.Request.FormFile("teams")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teams.csv file not provided"})
		return
	}
	defer teamsFile.Close()

	teams, err := parseTeamsCSV(teamsFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("parseTeamsCSV: %v", err.Error())})
		return
	}

	playersFile, _, err := c.Request.FormFile("players")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "players.csv file not provided"})
		return
	}
	defer playersFile.Close()

	players, err := parsePlayersCSV(playersFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("parsePlayersCSV: %v", err.Error())})
		return
	}

	if err := InsertTeamData(h.db, teams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertTeamData: %v", err)})
		return
	}
	if err := InsertPlayerData(h.db, players); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertPlayerData: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

func parseTeamsCSV(file io.Reader) ([]map[string]string, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	reader.TrimLeadingSpace = true

	var teams []map[string]string

	firstRow := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading teams.csv: %v", err)
		}

		if firstRow {
			firstRow = false
			continue
		}

		team := map[string]string{
			"name":      record[0],
			"captain":   record[1],
			"firstYear": record[2],
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func parsePlayersCSV(file io.Reader) ([]Player, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 7
	reader.TrimLeadingSpace = true

	var players []Player

	firstRow := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading players.csv: %v", err)
		}

		if firstRow {
			firstRow = false
			continue
		}

		player := Player{
			firstName: record[0],
			lastName:  record[1],
			email:     record[2],
			phone:     record[3],
			usaNum:    record[4],
			level:     record[5],
			teamNames: record[6],
		}
		players = append(players, player)
	}

	return players, nil
}
