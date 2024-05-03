package calendar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/api/calendar/v3"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/calendar/events", h.listEvents).Methods(http.MethodGet)
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	calendarID := "tpfitz42@gmail.com"
	currentYear := time.Now().Year()
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(currentYear, 12, 31, 23, 59, 59, 0, time.UTC)
	timeMin := startOfYear.Format(time.RFC3339)
	timeMax := endOfYear.Format(time.RFC3339)

	events, err := h.gcp.Events.List(calendarID).
		TimeMin(timeMin).
		TimeMax(timeMax).
		Context(h.ctx).
		Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve events: %v", err), http.StatusInternalServerError)
		return
	}

	calendars, err := h.gcp.CalendarList.List().Context(h.ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve calendar list: %v", err), http.StatusInternalServerError)
		return
	}

	type Response struct {
		Events    []*calendar.Event             `json:"events"`
		Calendars []*calendar.CalendarListEntry `json:"calendars"`
	}

	resp := Response{
		Events:    events.Items,
		Calendars: calendars.Items,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
