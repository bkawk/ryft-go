package ryft

import (
	"context"
	"net/url"
)

type DisputesService struct {
	client *Client
}

type Dispute struct {
	ID               string         `json:"id,omitempty"`
	Status           string         `json:"status,omitempty"`
	Category         string         `json:"category,omitempty"`
	CreatedTimestamp int            `json:"createdTimestamp,omitempty"`
	Reason           map[string]any `json:"reason,omitempty"`
	Files            map[string]any `json:"files,omitempty"`
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

func (s *DisputesService) List(
	ctx context.Context,
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	startsAfter string,
) (*DisputeList, error) {
	query := url.Values{}
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}
	query.Set("ascending", boolString(ascending))
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}
	if startsAfter != "" {
		query.Set("startsAfter", startsAfter)
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "disputes", query, nil)
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
	req, err := s.client.newRequest(ctx, "GET", "disputes/"+disputeID, nil)
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
	req, err := s.client.newRequest(ctx, "POST", "disputes/"+disputeID+"/accept", nil)
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
	req, err := s.client.newRequest(ctx, "POST", "disputes/"+disputeID+"/challenge", nil)
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
	req, err := s.client.newRequest(ctx, "PATCH", "disputes/"+disputeID+"/evidence", request)
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
	req, err := s.client.newRequest(ctx, "DELETE", "disputes/"+disputeID+"/evidence", request)
	if err != nil {
		return nil, err
	}

	var dispute Dispute
	if err := s.client.doJSON(req, &dispute); err != nil {
		return nil, err
	}
	return &dispute, nil
}
