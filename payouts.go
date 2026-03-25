package ryft

import (
	"context"
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

func (s *PayoutsService) Create(ctx context.Context, accountID string, request CreatePayoutRequest) (*Payout, error) {
	req, err := s.client.newRequest(ctx, "POST", "accounts/"+accountID+"/payouts", request)
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
	req, err := s.client.newRequest(ctx, "GET", "accounts/"+accountID+"/payouts/"+payoutID, nil)
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
	startTimestamp int,
	endTimestamp int,
	ascending bool,
	limit int,
	startsAfter string,
) (*PayoutList, error) {
	query := buildListQuery(ascending, limit, startsAfter)
	if startTimestamp > 0 {
		query.Set("startTimestamp", itoa(startTimestamp))
	}
	if endTimestamp > 0 {
		query.Set("endTimestamp", itoa(endTimestamp))
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "accounts/"+accountID+"/payouts", query, nil)
	if err != nil {
		return nil, err
	}

	var payouts PayoutList
	if err := s.client.doJSON(req, &payouts); err != nil {
		return nil, err
	}

	return &payouts, nil
}
