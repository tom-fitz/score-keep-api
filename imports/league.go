package imports

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) importLeagues(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Set max memory for file uploads (10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("ParseMultipartForm: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	teamsFile, _, err := r.FormFile("teams")
	if err != nil {
		http.Error(w, "teams.csv file not provided", http.StatusBadRequest)
		return
	}
	defer teamsFile.Close()

	teams, err := parseTeamsCSV(teamsFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("parseTeamsCSV: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	playersFile, _, err := r.FormFile("players")
	if err != nil {
		http.Error(w, "players.csv file not provided", http.StatusBadRequest)
		return
	}
	defer playersFile.Close()

	players, err := parsePlayersCSV(playersFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsePlayersCSV: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("players: %v", players)
	h.logger.Printf("teams: %v", teams)

	//resp := map[string]string{
	//	"status": "ok",
	//	"method": "post",
	//}
	//bytes, err := json.Marshal(resp)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(bytes)
}

func parseTeamsCSV(file io.Reader) ([]map[string]string, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	reader.TrimLeadingSpace = true

	var teams []map[string]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading teams.csv: %v", err)
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

func parsePlayersCSV(file io.Reader) ([]map[string]string, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 7
	reader.TrimLeadingSpace = true

	var players []map[string]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading players.csv: %v", err)
		}

		player := map[string]string{
			"firstName": record[0],
			"lastName":  record[1],
			"email":     record[2],
			"phone":     record[3],
			"usaNum":    record[4],
			"level":     record[5],
			"teamNames": record[6],
		}
		players = append(players, player)
	}

	return players, nil
}
