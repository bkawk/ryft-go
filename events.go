package ryft

import (
	"context"
	"encoding/json"
	"net/http"
)

type EventsService struct {
	client *Client
}

type Event struct {
	ID        string          `json:"id,omitempty"`
	EventType string          `json:"eventType,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	AccountID string          `json:"accountId,omitempty"`
}

type EventList struct {
	Items []Event `json:"items"`
}

type EventListParams struct {
	ListParams
}

func (s *EventsService) List(ctx context.Context, params EventListParams, opts ...RequestOption) (*EventList, error) {
	query := buildListQuery(params.ListParams)

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "events", query, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var events EventList
	if err := s.client.doJSON(req, &events); err != nil {
		return nil, err
	}
	return &events, nil
}

func (s *EventsService) Get(ctx context.Context, eventID string, opts ...RequestOption) (*Event, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "events/"+eventID, nil)
	if err != nil {
		return nil, err
	}
	s.client.applyRequestOptions(req, opts...)

	var event Event
	if err := s.client.doJSON(req, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
