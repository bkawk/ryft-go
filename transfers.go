package ryft

import (
	"context"
	"encoding/json"
	"net/http"
)

type TransfersService struct {
	client *Client
}

type Transfer struct {
	ID          string          `json:"id"`
	Amount      int             `json:"amount,omitempty"`
	Currency    string          `json:"currency,omitempty"`
	Status      string          `json:"status,omitempty"`
	Destination json.RawMessage `json:"destination,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

type TransferList struct {
	Items []Transfer `json:"items"`
}

type TransferListParams struct {
	ListParams
}

type TransferDestination struct {
	AccountID string `json:"accountId"`
}

type CreateTransferRequest struct {
	Amount      int                 `json:"amount"`
	Currency    string              `json:"currency"`
	Destination TransferDestination `json:"destination"`
	Reason      string              `json:"reason"`
	Metadata    map[string]string   `json:"metadata,omitempty"`
}

func (s *TransfersService) Create(ctx context.Context, request CreateTransferRequest) (*Transfer, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "transfers", request)
	if err != nil {
		return nil, err
	}

	var transfer Transfer
	if err := s.client.doJSON(req, &transfer); err != nil {
		return nil, err
	}

	return &transfer, nil
}

func (s *TransfersService) Get(ctx context.Context, transferID string) (*Transfer, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "transfers/"+transferID, nil)
	if err != nil {
		return nil, err
	}

	var transfer Transfer
	if err := s.client.doJSON(req, &transfer); err != nil {
		return nil, err
	}

	return &transfer, nil
}

func (s *TransfersService) List(ctx context.Context, params TransferListParams) (*TransferList, error) {
	query := buildListQuery(params.ListParams)

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "transfers", query, nil)
	if err != nil {
		return nil, err
	}

	var transfers TransferList
	if err := s.client.doJSON(req, &transfers); err != nil {
		return nil, err
	}

	return &transfers, nil
}
