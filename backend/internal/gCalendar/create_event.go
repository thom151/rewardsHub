package gCalendar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CreateEventParams struct {
	CalendarID  string
	Summary     string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	TimeZone    string
}

type EventResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	HtmlLink    string `json:"htmlLink"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Start       struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
}

func CreateEvent(ctx context.Context, accToken string, p CreateEventParams) (EventResponse, error) {
	if p.CalendarID == "" {
		p.CalendarID = "primary"
	}
	if p.TimeZone == "" {
		p.TimeZone = "Australia/Melbourne"
	}

	url := fmt.Sprintf(
		"https://www.googleapis.com/calendar/v3/calendars/%s/events",
		p.CalendarID,
	)

	bodyMap := map[string]any{
		"summary":     p.Summary,
		"description": p.Description,
		"location":    p.Location,
		"start": map[string]any{
			"dateTime": p.Start.Format(time.RFC3339),
			"timeZone": p.TimeZone,
		},
		"end": map[string]any{
			"dateTime": p.End.Format(time.RFC3339),
			"timeZone": p.TimeZone,
		},
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return EventResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return EventResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EventResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EventResponse{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return EventResponse{}, fmt.Errorf("google calendar create event error (%d): %s", resp.StatusCode, string(body))
	}

	var event EventResponse
	if err = json.Unmarshal(body, &event); err != nil {
		return EventResponse{}, err
	}

	return event, nil

}
