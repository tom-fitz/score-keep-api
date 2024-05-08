package imports

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Team struct {
	Name      string `json:"name"`
	Captain   string `json:"captain"`
	FirstYear string `json:"firstYear"`
}

func (h *Service) ImportTeams(c *gin.Context) {
	lid := c.Param("leagueId")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	if err := validateLeague(lid, h.db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
		return
	}

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

	if err := InsertTeamData(h.db, teams); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertTeamData: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teams imported successfully"})
}

func parseTeamsCSV(file io.Reader) ([]Team, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var teams []Team

	firstRow := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading teams CSV: %v", err)
		}

		if firstRow {
			firstRow = false
			continue
		}

		team, err := parseTeamRecord(record)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func parseTeamRecord(record []string) (Team, error) {
	// TODO: add in team validation
	team := Team{
		Name:      record[0],
		Captain:   record[1],
		FirstYear: record[2],
	}
	return team, nil
}
