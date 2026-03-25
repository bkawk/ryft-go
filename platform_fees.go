package ryft

import (
	"context"
	"net/http"
)

type PlatformFeesService struct {
	client *Client
}

type PlatformFee struct {
	ID               string `json:"id,omitempty"`
	PaymentSessionID string `json:"paymentSessionId,omitempty"`
	FromAccountID    string `json:"fromAccountId,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	Currency         string `json:"currency,omitempty"`
}

type PlatformFeeList struct {
	Items []PlatformFee `json:"items"`
}

type PlatformFeeRefund struct {
	ID            string `json:"id,omitempty"`
	PlatformFeeID string `json:"platformFeeId,omitempty"`
	Amount        int    `json:"amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	Status        string `json:"status,omitempty"`
}

type PlatformFeeRefundList struct {
	Items []PlatformFeeRefund `json:"items"`
}

type PlatformFeeListParams struct {
	ListParams
}

func (s *PlatformFeesService) List(ctx context.Context, params PlatformFeeListParams) (*PlatformFeeList, error) {
	query := buildListQuery(params.ListParams)

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "platform-fees", query, nil)
	if err != nil {
		return nil, err
	}

	var fees PlatformFeeList
	if err := s.client.doJSON(req, &fees); err != nil {
		return nil, err
	}
	return &fees, nil
}

func (s *PlatformFeesService) Get(ctx context.Context, feeID string) (*PlatformFee, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "platform-fees/"+feeID, nil)
	if err != nil {
		return nil, err
	}

	var fee PlatformFee
	if err := s.client.doJSON(req, &fee); err != nil {
		return nil, err
	}
	return &fee, nil
}

func (s *PlatformFeesService) GetRefunds(ctx context.Context, feeID string) (*PlatformFeeRefundList, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "platform-fees/"+feeID+"/refunds", nil)
	if err != nil {
		return nil, err
	}

	var refunds PlatformFeeRefundList
	if err := s.client.doJSON(req, &refunds); err != nil {
		return nil, err
	}
	return &refunds, nil
}
