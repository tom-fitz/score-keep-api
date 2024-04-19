package calendar

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func NewCalendarService(ctx context.Context, filename string) (*calendar.Service, error) {
	scopes := []string{calendar.CalendarScope, calendar.CalendarEventsScope}

	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	creds, err := google.CredentialsFromJSON(ctx, jsonBytes, scopes...)
	if err != nil {
		return nil, err
	}

	svc, err := calendar.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return svc, nil
}
