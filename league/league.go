package league

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) createLeagues(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	h.logger.Printf("Received POST request with JSON body: %v", requestBody)
}
