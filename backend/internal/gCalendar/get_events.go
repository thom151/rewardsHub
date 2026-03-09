package gCalendar

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type EventsResponse struct {
	Items []struct {
		Summary string `json:"summary"`
		Start   struct {
			DateTime string `json:"dateTime"`
			Date     string `json:"date"`
		} `json:"start"`
		End struct {
			DateTime string `json:"dateTime"`
			Date     string `json:"date"`
		} `json:"end"`
	} `json:"items"`
}

func GetCalendarEvents(ctx context.Context, accToken string, start, end time.Time) (EventsResponse, error) {
	base := "https://www.googleapis.com/calendar/v3/calendars/primary/events"

	params := url.Values{}
	params.Set("timeMin", start.UTC().Format(time.RFC3339))
	params.Set("timeMax", end.UTC().Format(time.RFC3339))
	params.Set("singleEvents", "true")
	params.Set("orderBy", "startTime")

	fullUrl := base + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
	if err != nil {
		return EventsResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+accToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EventsResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EventsResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return EventsResponse{}, err
	}

	var events EventsResponse
	if err := json.Unmarshal(body, &events); err != nil {
		return EventsResponse{}, err
	}

	return events, nil
}
