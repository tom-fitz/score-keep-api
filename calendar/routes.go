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
	//settings, err := h.gcp.Settings.List().Do()
	// https://calendar.google.com/calendar/u/0?cid=dHBmaXR6NDJAZ21haWwuY29t
	acl, err := h.gcp.Acl.List("primary").Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Printf("acl: %+v", acl)
	for _, item := range acl.Items {
		h.logger.Printf("acl 2: %+v", *item.Scope)
	}

	//tokenInfo, err := h.gcp.Tokeninfo().Context(h.ctx).Do()
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Unable to get token info: %v", err), http.StatusInternalServerError)
	//	return
	//}
	// Fetch events from the specified calendar
	events, err := h.gcp.Events.List("primary").Context(h.ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve events: %v", err), http.StatusInternalServerError)
		return
	}

	// ID of the calendar to insert
	calendarID := "dHBmaXR6NDJAZ21haWwuY29t"

	// Create a new CalendarListEntry object with the specified ID
	newCalendar := &calendar.CalendarListEntry{
		Id: calendarID,
	}

	calListEntry, err := h.gcp.CalendarList.Insert(newCalendar).Context(h.ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to insert calendar id: %s err: %v", "dHBmaXR6NDJAZ21haWwuY29t", err), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("calListEntry: %+v", calListEntry)

	// Fetch a list of all calendars the user has access to
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
