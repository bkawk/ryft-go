package ryft

import (
	"context"
	"net/http"
	"strconv"
)

type PayoutsService struct {
	client *Client
}

type Payout struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
	Status   string `json:"status,omitempty"`
}

type CreatePayoutRequest struct {
	Amount         int               `json:"amount"`
	Currency       string            `json:"currency"`
	PayoutMethodID string            `json:"payoutMethodId"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type PayoutList struct {
	Items []Payout `json:"items"`
}

type PayoutListParams struct {
	ListParams
	StartTimestamp int
	EndTimestamp   int
}

func (s *PayoutsService) Create(ctx context.Context, accountID string, request CreatePayoutRequest) (*Payout, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "accounts/"+accountID+"/payouts", request)
	if err != nil {
		return nil, err
	}

	var payout Payout
	if err := s.client.doJSON(req, &payout); err != nil {
		return nil, err
	}

	return &payout, nil
}

func (s *PayoutsService) Get(ctx context.Context, accountID string, payoutID string) (*Payout, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "accounts/"+accountID+"/payouts/"+payoutID, nil)
	if err != nil {
		return nil, err
	}

	var payout Payout
	if err := s.client.doJSON(req, &payout); err != nil {
		return nil, err
	}

	return &payout, nil
}

func (s *PayoutsService) List(
	ctx context.Context,
	accountID string,
	params PayoutListParams,
) (*PayoutList, error) {
	query := buildListQuery(params.ListParams)
	if params.StartTimestamp > 0 {
		query.Set("startTimestamp", strconv.Itoa(params.StartTimestamp))
	}
	if params.EndTimestamp > 0 {
		query.Set("endTimestamp", strconv.Itoa(params.EndTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "accounts/"+accountID+"/payouts", query, nil)
	if err != nil {
		return nil, err
	}

	var payouts PayoutList
	if err := s.client.doJSON(req, &payouts); err != nil {
		return nil, err
	}

	return &payouts, nil
}
