// league.go
package imports

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (h *Handler) importLeagues(c *gin.Context) {
	//id := c.Param("id")
	//if id == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "league id required"})
	//	return
	//}
	//
	//if err := validateLeague(id, h.db); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error validating league: %v", err)})
	//	return
	//}
	//
	//teamsFile, _, err := c.Request.FormFile("teams")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "teams.csv file not provided"})
	//	return
	//}
	//defer teamsFile.Close()
	//
	//teams, err := parseTeamCSV(teamsFile, parseTeamRecord)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("parseTeamsCSV: %v", err.Error())})
	//	return
	//}
	//
	//playersFile, _, err := c.Request.FormFile("players")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "players.csv file not provided"})
	//	return
	//}
	//defer playersFile.Close()
	//
	//players, err := parseCSV(playersFile, parsePlayerRecord)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("parsePlayersCSV: %v", err.Error())})
	//	return
	//}
	//
	//if err := InsertTeamData(h.db, teams); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertTeamData: %v", err)})
	//	return
	//}
	//if err := InsertPlayerData(h.db, players); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("InsertPlayerData: %v", err)})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

func parseCSV(file io.Reader, parseRecord func(record []string) (interface{}, error)) (interface{}, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var data interface{}

	firstRow := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}

		if firstRow {
			firstRow = false
			continue
		}

		item, err := parseRecord(record)
		if err != nil {
			return nil, err
		}
		data = append(data.([]interface{}), item)
	}

	return data, nil
}

//
//func parseTeamRecord(record []string) (Team, error) {
//	team := Team{
//		Name:      record[0],
//		Captain:   record[1],
//		FirstYear: record[2],
//	}
//	return team, nil
//}

func parsePlayerRecord(record []string) (Player, error) {
	player := Player{
		FirstName: record[0],
		LastName:  record[1],
		Email:     record[2],
		Phone:     record[3],
		UsaNum:    record[4],
		Level:     record[5],
		TeamNames: record[6],
	}
	return player, nil
}
