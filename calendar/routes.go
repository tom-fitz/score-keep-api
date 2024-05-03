package calendar

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/v1/calendar/events", h.listEvents)
}

func (h *Handler) listEvents(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to retrieve events: %v", err)})
		return
	}

	calendars, err := h.gcp.CalendarList.List().Context(h.ctx).Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to retrieve calendar list: %v", err)})
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

	c.JSON(http.StatusOK, resp)
}
