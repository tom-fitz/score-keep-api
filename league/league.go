package league

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) createLeagues(w http.ResponseWriter, r *http.Request) {
	var requestBody []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	var ids []int
	query := "INSERT INTO score_keep_db.public.leagues (name, level) VALUES ($1, $2) RETURNING id"
	for _, league := range requestBody {
		var id int
		err := h.db.QueryRow(query, league["name"], league["level"]).Scan(&id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to insert data into database: %v", err), http.StatusInternalServerError)
			return
		}
		ids = append(ids, id)
	}

	json.NewEncoder(w).Encode(map[string][]int{"ids": ids})
}
