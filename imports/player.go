package imports

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Player struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Usanum    string `json:"usanum"`
	Level     string `json:"level"`
	TeamNames string `json:"teamNames"`
}

func (h *Service) ImportPlayers(c *gin.Context) {
	lid := c.Param("leagueId")
	if lid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
		return
	}

	_, err := validateLeague(lid, h.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
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

	h.logger.Printf("parsedPlayers: %v", players)

	if err := InsertPlayerData(h.db, players); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertPlayerData: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Players imported successfully."})
}

func parsePlayersCSV(file io.Reader) ([]Player, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var players []Player

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

		player, err := parsePlayerRecord(record)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func parsePlayerRecord(record []string) (Player, error) {
	// TODO: add in team validation

	player := Player{
		FirstName: record[0],
		LastName:  record[1],
		Email:     record[2],
		Phone:     record[3],
		Usanum:    record[4],
		Level:     record[5],
		TeamNames: record[6],
	}
	return player, nil
}
