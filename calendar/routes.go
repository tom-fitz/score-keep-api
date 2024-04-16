package calendar

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/calendar/events", h.listEvents).Methods(http.MethodGet)
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.gcp.Events.List("primary").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding events to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write events JSON to response
	_, err = w.Write(eventsJSON)
	if err != nil {
		h.logger.Printf("Error writing events JSON to response: %v", err)
	}
}
