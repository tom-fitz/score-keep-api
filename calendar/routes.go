package calendar

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/calendar/v3"
	"net/http"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/calendar/events", h.listEvents).Methods(http.MethodGet)
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	// ID of the calendar to access
	calendarID := "tpfitz42@gmail.com"

	// Fetch events from the specified calendar
	events, err := h.gcp.Events.List(calendarID).Context(h.ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve events: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch a list of all calendars the service account has access to
	calendars, err := h.gcp.CalendarList.List().Context(h.ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar list: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a struct to hold both events and calendars data
	type Response struct {
		Events    []*calendar.Event             `json:"events"`
		Calendars []*calendar.CalendarListEntry `json:"calendars"`
	}

	// Create the response object
	resp := Response{
		Events:    events.Items,
		Calendars: calendars.Items,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write response JSON to response
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
