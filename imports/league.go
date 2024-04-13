package imports

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) importLeagues(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("test post...")
	resp := map[string]string{
		"status": "ok",
		"method": "post",
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}
