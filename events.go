package ryft

import (
	"context"
	"net/url"
)

type EventsService struct {
	client *Client
}

type Event struct {
	ID        string         `json:"id,omitempty"`
	EventType string         `json:"eventType,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	AccountID string         `json:"accountId,omitempty"`
}

type EventList struct {
	Items []Event `json:"items"`
}

func (s *EventsService) List(ctx context.Context, ascending bool, limit int, accountID string) (*EventList, error) {
	query := url.Values{}
	query.Set("ascending", boolString(ascending))
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "events", query, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var events EventList
	if err := s.client.doJSON(req, &events); err != nil {
		return nil, err
	}
	return &events, nil
}

func (s *EventsService) Get(ctx context.Context, eventID string, accountID string) (*Event, error) {
	req, err := s.client.newRequest(ctx, "GET", "events/"+eventID, nil)
	if err != nil {
		return nil, err
	}
	if accountID != "" {
		req.Header.Set("Account", accountID)
	}

	var event Event
	if err := s.client.doJSON(req, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
