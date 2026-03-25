package ryft

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type DisputesService struct {
	client *Client
}

type Dispute struct {
	ID               string          `json:"id,omitempty"`
	Status           string          `json:"status,omitempty"`
	Category         string          `json:"category,omitempty"`
	CreatedTimestamp int             `json:"createdTimestamp,omitempty"`
	Reason           json.RawMessage `json:"reason,omitempty"`
	Files            json.RawMessage `json:"files,omitempty"`
}

type DisputeList struct {
	Items []Dispute `json:"items"`
}

type AddDisputeEvidenceRequest struct {
	Files map[string]any `json:"files,omitempty"`
	Text  map[string]any `json:"text,omitempty"`
}

type DeleteDisputeEvidenceRequest struct {
	Files []string `json:"files"`
}

type DisputeListParams struct {
	ListParams
	StartTimestamp int
	EndTimestamp   int
}

func (s *DisputesService) List(ctx context.Context, params DisputeListParams) (*DisputeList, error) {
	query := buildListQuery(params.ListParams)
	if params.StartTimestamp > 0 {
		query.Set("startTimestamp", strconv.Itoa(params.StartTimestamp))
	}
	if params.EndTimestamp > 0 {
		query.Set("endTimestamp", strconv.Itoa(params.EndTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "disputes", query, nil)
	if err != nil {
		return nil, err
	}

	var disputes DisputeList
	if err := s.client.doJSON(req, &disputes); err != nil {
		return nil, err
	}
	return &disputes, nil
}

func (s *DisputesService) Get(ctx context.Context, disputeID string) (*Dispute, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "disputes/"+disputeID, nil)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (s *DisputesService) Accept(ctx context.Context, disputeID string) (*Dispute, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "disputes/"+disputeID+"/accept", nil)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (s *DisputesService) Challenge(ctx context.Context, disputeID string) (*Dispute, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "disputes/"+disputeID+"/challenge", nil)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (s *DisputesService) AddEvidence(ctx context.Context, disputeID string, request AddDisputeEvidenceRequest) (*Dispute, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, "disputes/"+disputeID+"/evidence", request)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}

func (s *DisputesService) DeleteEvidence(ctx context.Context, disputeID string, request DeleteDisputeEvidenceRequest) (*Dispute, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, "disputes/"+disputeID+"/evidence", request)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}
